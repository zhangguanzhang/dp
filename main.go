package main

import (
	"docker-pull/cmd"
)


func main() {
	const (
		ImgFmt = ".docker_temp_%s"
	)

	//p := registry.NewPull(os.Args[1])
	//data, err := p.Manifests()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(p.Repository)
	//fmt.Println(data.Layers[0].Descriptor())
	cmd.Execute()

	//log.Println(registry.Save(os.Args[1:], "test.tar.gz"))



	//imgDir := fmt.Sprintf(ImgFmt, strings.ReplaceAll(os.Args[1], ":", "@"))
	//// We use sequential file access here to avoid depleting the standby list
	//// on Windows. On Linux, this is a call directly to ioutil.TempFile
	//tmpFile, err := pkg.TempFileSequential(filepath.Dir(imgDir), ".docker_temp_")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//tmpPath := tmpFile.Name()

	//_, Bytes, err := p.Blobs(data.Layers[0].Digest)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(io.Copy(ioutil.Discard, bytes.NewReader(Bytes)))
}
