package registry

import "time"

//https://github.com/moby/moby/blob/master/image/image.go
//https://github.com/moby/moby/blob/master/api/types/container/config.go


type ImageConfig struct {
	Architecture string `json:"architecture"`
	Config       struct {
		AttachStderr bool        `json:"AttachStderr"`
		AttachStdin  bool        `json:"AttachStdin"`
		AttachStdout bool        `json:"AttachStdout"`
		OpenStdin  bool          `json:"OpenStdin"`
		StdinOnce  bool          `json:"StdinOnce"`
		Tty        bool          `json:"Tty"`
		Cmd          []string    `json:"Cmd"`
		Domainname   string      `json:"Domainname"`
		Entrypoint   interface{} `json:"Entrypoint"`
		Env          []string    `json:"Env"`
		Hostname     string      `json:"Hostname"`
		Image        string      `json:"Image"`
		Labels       struct {
		} `json:"Labels"`
		OnBuild    []interface{} `json:"OnBuild"`
		User       string        `json:"User"`
		Volumes    interface{}   `json:"Volumes"`
		WorkingDir string        `json:"WorkingDir"`
	} `json:"config"`
	Container       string `json:"container"`
	ContainerConfig struct {
		AttachStderr bool        `json:"AttachStderr"`
		AttachStdin  bool        `json:"AttachStdin"`
		AttachStdout bool        `json:"AttachStdout"`
		OpenStdin  bool          `json:"OpenStdin"`
		StdinOnce  bool          `json:"StdinOnce"`
		Tty        bool          `json:"Tty"`
		Cmd          []string    `json:"Cmd"`
		Domainname   string      `json:"Domainname"`
		Entrypoint   interface{} `json:"Entrypoint"`
		Env          []string    `json:"Env"`
		Hostname     string      `json:"Hostname"`
		Image        string      `json:"Image"`
		Labels       struct {
		} `json:"Labels"`
		OnBuild    []interface{} `json:"OnBuild"`
		User       string        `json:"User"`
		Volumes    interface{}   `json:"Volumes"`
		WorkingDir string        `json:"WorkingDir"`
	} `json:"container_config"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	Os     string `json:"os"`
	//History       []struct {
	//	Created    time.Time `json:"created"`
	//	CreatedBy  string    `json:"created_by"`
	//	EmptyLayer bool      `json:"empty_layer,omitempty"`
	//} `json:"history"`
	//Rootfs struct {
	//	DiffIds []string `json:"diff_ids"`
	//	Type    string   `json:"type"`
	//} `json:"rootfs"`
	ID string `json:"id"`
	Parent string `json:"parent"`
}