package v1alpha1

type GenericDependencySpec struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type DependenciesSpec struct {
	OS     []GenericDependencySpec `json:"os,omitempty" description:"Not yet implemented"`
	Pip    []GenericDependencySpec `json:"pip,omitempty" description:"Not yet implemented"`
	Go     []GenericDependencySpec `json:"go,omitempty" description:"Not yet implemented"`
	Cargo  []GenericDependencySpec `json:"cargo,omitempty" description:"Not yet implemented"`
	NPM    []GenericDependencySpec `json:"npm,omitempty" description:"Not yet implemented"`
	Yarn   []GenericDependencySpec `json:"yarn,omitempty" description:"Not yet implemented"`
	ASDF   []GenericDependencySpec `json:"asdf,omitempty" description:"Not yet implemented"`
	Maven  []GenericDependencySpec `json:"maven,omitempty" description:"Not yet implemented"`
	Gradle []GenericDependencySpec `json:"gradle,omitempty" description:"Not yet implemented"`
}

type SourceCodeSpec struct {
	URL           string           `json:"url" description:"URL of the git repo of source code"`
	Branch        string           `json:"branch,omitempty"`
	Commit        string           `json:"commit,omitempty"`
	Dependencies  DependenciesSpec `json:"dependencies,omitempty"`
	DependencyCMD string           `json:"dependencycmd,omitempty"`
	BuildCMD      string           `json:"buildcmd,omitempty"`
	Entrypoint    string           `json:"entrypoint"`
}

type InputDataSpec struct {
	URL          string `json:"url" description:"URL of the data repo of input data"`
	Type         string `json:"type" description:"Type of the input data source, e.g. s3, git, http, https, etc."`
	DataPath     string `json:"datapath" description:"Path to the data directory with input data, has to be a unix path."`
	TransformCMD string `json:"transformcmd,omitempty" description:"Command to transform the input data"`
}

type OutputDataSpec struct {
	URL      string `json:"url" description:"URL of the data repo of output data"`
	DataPath string `json:"datapath" description:"Path to the data directory with output data, has to be a unix path."`
}
