package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/opencontainers/go-digest"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"strings"
)

const  (
	REPO = "library"
	DefaultTAG = "latest"
	REGISTRY =  "registry-1.docker.io"
	V2JSON = "application/vnd.docker.distribution.manifest.v2+json"
	//github.com/docker/docker/image/tarexport
	ManifestFileName           = "manifest.json"
	LegacyLayerFileName        = "layer.tar"
	LegacyConfigFileName       = "json"
	LegacyVersionFileName      = "VERSION"
	LegacyRepositoriesFileName = "repositories"
)

type WriteBarFunc func(downloadName string, length, downLen int64)

type TarAddfileFunc func(size int64, name string, b interface{}) error

type Pull struct {
	// like registry-1.docker.io
	Registry string
	//like latest
	Tag string
	// Registry namespcase
	Repository string
	// Must be implemented in order to verify `RoundTrip()`
	Client *http.Client
	//
	ImgParts []string
}

// for manifest.json file
type ManifestItem struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

func NewPull(pullImg string) *Pull {
	p := &Pull{Tag: DefaultTAG}
	repo := REPO
	tempStrSlice := make([]string, 0)
	imgParts := strings.Split(pullImg, "/")
	if strings.Contains(imgParts[len(imgParts)-1], "@") {
		tempStrSlice = strings.Split(imgParts[len(imgParts)-1], "@")
	} else if strings.Contains(imgParts[len(imgParts)-1], ":"){
		tempStrSlice = strings.Split(imgParts[len(imgParts)-1], ":")
	} else {
		tempStrSlice = []string{imgParts[len(imgParts)-1], DefaultTAG}
	}
	img := tempStrSlice[0]
	p.Tag = tempStrSlice[1]

	//`:` means the port, the first part has `.` means the domain name or ip
	if len(imgParts) > 1 &&
		( strings.Contains(imgParts[0], ".") || strings.Contains(imgParts[0], ":") ) {
		// use domain
		p.Registry = imgParts[0]
		repo = strings.Join(imgParts[1:len(imgParts) - 1], "/")
	} else {// dockerhub
		p.Registry = REGISTRY
		if len(imgParts[:len(imgParts)-1]) != 0 {
			repo = strings.Join(imgParts[:len(imgParts)-1], "/")
		}
	}
	p.Repository = fmt.Sprintf("%s/%s", repo, img)

	p.Client = &http.Client{
		Transport: NewTokenTransport(&http.Transport{
			Proxy:                  http.ProxyFromEnvironment,
			//DialContext:  (&net.Dialer{
			//	Timeout:   30 * time.Second,
			//	KeepAlive: 30 * time.Second,
			//}).DialContext,
			//ForceAttemptHTTP2:     true,
			//MaxIdleConns:          150,
			//MaxIdleConnsPerHost:   -1,
			//IdleConnTimeout:       90 * time.Second,
			//TLSHandshakeTimeout:   10 * time.Second,
			//ExpectContinueTimeout: 5 * time.Second,
			TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
			//ResponseHeaderTimeout: time.Second * 8,
		}),
		Timeout:   time.Second * 15,
	}

	return p
}

func (p *Pull) Do(req *http.Request) (*http.Response, error) {
	return p.Client.Do(req)
}

func (p *Pull) Manifests() (*schema2.Manifest, error) {
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://%s/v2/%s/manifests/%s", p.Registry, p.Repository, p.Tag), nil)
	req.Header.Set("Accept", V2JSON)
	resp, err := p.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while request manifests|%s", err)
	}
	defer resp.Body.Close()

	respBody,_ := ioutil.ReadAll(resp.Body)

	var data schema2.Manifest
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, fmt.Errorf("unmarshal err|%s", err)
	}
	return &data, nil
}


func (p *Pull) Blobs(Digest digest.Digest) (int64, io.ReadCloser, error) {
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://%s/v2/%s/blobs/%s", p.Registry, p.Repository, Digest.String()), nil)
	req.Header.Set("Accept", V2JSON)
	resp, err := p.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("while request blobs|%s", err)
	}

	fSize, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("Content-Length|%s", err)
	}
	return fSize, resp.Body, nil
}


func Save(names []string, fileName string) (error) {
	fw, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	tarAddfile := TarAddfileWithDownBar(tw, WriteBar)
	var (
		manifestJsons = make([]ManifestItem, 0)
		repositoriesJson = make(map[string]map[string]string, 1)
	)


	for _, name := range names {
		parentID := ""
		p := NewPull(name)
		data, err := p.Manifests()
		if err != nil {
			return err
		}

		fSize, confRespsBody, err := p.Blobs(data.Config.Digest)
		if err != nil {
			return err
		}
		// for id.json file
		confResps, err := ioutil.ReadAll(confRespsBody)
		defer confRespsBody.Close()
		if err := tarAddfile(fSize, data.Config.Digest.Hex() + ".json", confResps);err != nil {
			return err
		}
		manifestJson := ManifestItem{
			Config: data.Config.Digest.Hex() + ".json",
			RepoTags: []string{name},
			Layers: make([]string, 0),
		}
		var confCont ImageConfig
		err = json.Unmarshal(confResps, &confCont)
		if err != nil {
			return err
		}

		for i, layer := range data.Layers {
			_, respBody, err := p.Blobs(layer.Digest)
			if err != nil {
				return err
			}

			// https://github.com/moby/moby/blob/master/image/tarexport/save.go#L294-L329
			// https://gist.github.com/aaronlehmann/b42a2eaf633fc949f93b#id-definitions-and-calculations
			legacyLayerDir := fmt.Sprintf("%x",
				sha256.Sum256([]byte(fmt.Sprintf(parentID + "\n" + layer.Digest.String() + "\n"))))

			// for layer.tar
			if err := tarAddfile(layer.Size, filepath.Join(legacyLayerDir, LegacyLayerFileName), respBody);err != nil {
				return err
			}
			manifestJson.Layers = append(manifestJson.Layers, filepath.Join(legacyLayerDir, LegacyLayerFileName))

			// for VERSION
			if err := tarAddfile(int64(len([]byte(`1.0`))),
					filepath.Join(legacyLayerDir, LegacyVersionFileName), []byte(`1.0`));err != nil {
				return err
			}

			confCont.ID = legacyLayerDir

			if parentID != "" {
				confCont.Parent = parentID
			}
			parentID = confCont.ID

			confBytes := NewLayerEmptyJson()
			if i == len(data.Layers) - 1 {
				confBytesFull, _:= json.Marshal(&confCont)
				confBytes = confBytesFull
			}

			// for json
			if err := tarAddfile(int64(len(confBytes)),
				filepath.Join(legacyLayerDir, LegacyConfigFileName), confBytes);err != nil {
				return err
			}
			//layerName append
		}
		manifestJsons = append(manifestJsons, manifestJson)
		if v, ok := repositoriesJson[p.Repository]; ok {
			v[p.Tag] = data.Layers[len(data.Layers)-1].Digest.Hex()
		} else {
			repositoriesJson[p.Repository] = map[string]string{
				p.Tag: data.Layers[len(data.Layers)-1].Digest.Hex(),
			}
		}
	}

	//for ManifestFileName   = "manifest.json"
	manifestBytes, err := json.Marshal(&manifestJsons)
	if err != nil {
		return fmt.Errorf("while Marshal manifestJsons|%s", err)
	}
	if err := tarAddfile(int64(len(manifestBytes)), ManifestFileName, manifestBytes);err != nil {
		return err
	}

	//for LegacyRepositoriesFileName = "repositories"
	repositoriesBytes, err := json.Marshal(&repositoriesJson)
	if err != nil {
		return fmt.Errorf("while Marshal repositoriesBytes|%s", err)
	}
	if err := tarAddfile(int64(len(repositoriesBytes)), LegacyRepositoriesFileName, repositoriesBytes);err != nil {
		return err
	}

	return nil
}

func WriteBar(downloadName string, length, downLen int64) {
	fmt.Printf("\r%-76s CurrentTotalBytes %15d, ConsumedTotalBytes: %15d, %d%%",
		downloadName, length, downLen, downLen*100/length)
}

func TarAddfileWithDownBar(tw *tar.Writer, wb WriteBarFunc) TarAddfileFunc {
	return func(size int64, name string, b interface{}) error {
		var (
			buf     = make([]byte, 32*1024)
			written int64
			err error
			data  io.Reader
		)

		err = tw.WriteHeader(&tar.Header{
			Mode: 0644,
			Size: size,
			Name: name,
		})
		if err != nil {
			return fmt.Errorf("%s write header|%s", name, err)
		}

		switch v := b.(type) {
		case []byte:
			data = bytes.NewReader(v)
		case io.ReadCloser:
			data = v
		default:
			return fmt.Errorf("invalid type")
		}
		
		for {
			numRead, readErr := data.Read(buf)
			if numRead > 0 {
				numWrite, writeErr := tw.Write(buf[0:numRead])
				if numWrite > 0 {
					written += int64(numWrite)
				}
				if writeErr != nil {
					err = io.ErrShortWrite
					break
				}
			}
			if readErr != nil {
				if readErr != io.EOF {
					err = readErr
				}
				break
			}
			wb(name, size, written)
		}
		fmt.Println()
		//_, err = tw.Write(b)
		//if err != nil {
		//	return fmt.Errorf("%s write bytes|%s", name, err)
		//}
		//if err := tw.Flush();err != nil {
		//	return err
		//}
		return nil
	}
}

type EmptyConfig struct {
	Created         time.Time `json:"created"`
	ContainerConfig struct {
		Hostname     string      `json:"Hostname"`
		Domainname   string      `json:"Domainname"`
		User         string      `json:"User"`
		AttachStdin  bool        `json:"AttachStdin"`
		AttachStdout bool        `json:"AttachStdout"`
		AttachStderr bool        `json:"AttachStderr"`
		Tty          bool        `json:"Tty"`
		OpenStdin    bool        `json:"OpenStdin"`
		StdinOnce    bool        `json:"StdinOnce"`
		Env          interface{} `json:"Env"`
		Cmd          interface{} `json:"Cmd"`
		Image        string      `json:"Image"`
		Volumes      interface{} `json:"Volumes"`
		WorkingDir   string      `json:"WorkingDir"`
		Entrypoint   interface{} `json:"Entrypoint"`
		OnBuild      interface{} `json:"OnBuild"`
		Labels       interface{} `json:"Labels"`
	} `json:"container_config"`
}

func NewLayerEmptyJson() []byte {

	d, _ := json.Marshal(&EmptyConfig{
		Created: time.Unix(0, 0),
		ContainerConfig: struct {
			Hostname     string      `json:"Hostname"`
			Domainname   string      `json:"Domainname"`
			User         string      `json:"User"`
			AttachStdin  bool        `json:"AttachStdin"`
			AttachStdout bool        `json:"AttachStdout"`
			AttachStderr bool        `json:"AttachStderr"`
			Tty          bool        `json:"Tty"`
			OpenStdin    bool        `json:"OpenStdin"`
			StdinOnce    bool        `json:"StdinOnce"`
			Env          interface{} `json:"Env"`
			Cmd          interface{} `json:"Cmd"`
			Image        string      `json:"Image"`
			Volumes      interface{} `json:"Volumes"`
			WorkingDir   string      `json:"WorkingDir"`
			Entrypoint   interface{} `json:"Entrypoint"`
			OnBuild      interface{} `json:"OnBuild"`
			Labels       interface{} `json:"Labels"`
		}{},
	})
	return d
}
