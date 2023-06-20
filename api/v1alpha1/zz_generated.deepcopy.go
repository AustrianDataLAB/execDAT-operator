//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 Thomas Weber.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Build) DeepCopyInto(out *Build) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Build.
func (in *Build) DeepCopy() *Build {
	if in == nil {
		return nil
	}
	out := new(Build)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Build) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuildList) DeepCopyInto(out *BuildList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Build, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuildList.
func (in *BuildList) DeepCopy() *BuildList {
	if in == nil {
		return nil
	}
	out := new(BuildList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BuildList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuildSpec) DeepCopyInto(out *BuildSpec) {
	*out = *in
	in.SourceCode.DeepCopyInto(&out.SourceCode)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuildSpec.
func (in *BuildSpec) DeepCopy() *BuildSpec {
	if in == nil {
		return nil
	}
	out := new(BuildSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BuildStatus) DeepCopyInto(out *BuildStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BuildStatus.
func (in *BuildStatus) DeepCopy() *BuildStatus {
	if in == nil {
		return nil
	}
	out := new(BuildStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependenciesSpec) DeepCopyInto(out *DependenciesSpec) {
	*out = *in
	if in.OS != nil {
		in, out := &in.OS, &out.OS
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Pip != nil {
		in, out := &in.Pip, &out.Pip
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Go != nil {
		in, out := &in.Go, &out.Go
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Cargo != nil {
		in, out := &in.Cargo, &out.Cargo
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.NPM != nil {
		in, out := &in.NPM, &out.NPM
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Yarn != nil {
		in, out := &in.Yarn, &out.Yarn
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.ASDF != nil {
		in, out := &in.ASDF, &out.ASDF
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Maven != nil {
		in, out := &in.Maven, &out.Maven
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
	if in.Gradle != nil {
		in, out := &in.Gradle, &out.Gradle
		*out = make([]GenericDependencySpec, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependenciesSpec.
func (in *DependenciesSpec) DeepCopy() *DependenciesSpec {
	if in == nil {
		return nil
	}
	out := new(DependenciesSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericDependencySpec) DeepCopyInto(out *GenericDependencySpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericDependencySpec.
func (in *GenericDependencySpec) DeepCopy() *GenericDependencySpec {
	if in == nil {
		return nil
	}
	out := new(GenericDependencySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InputDataSpec) DeepCopyInto(out *InputDataSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InputDataSpec.
func (in *InputDataSpec) DeepCopy() *InputDataSpec {
	if in == nil {
		return nil
	}
	out := new(InputDataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutputDataSpec) DeepCopyInto(out *OutputDataSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutputDataSpec.
func (in *OutputDataSpec) DeepCopy() *OutputDataSpec {
	if in == nil {
		return nil
	}
	out := new(OutputDataSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodSpecData) DeepCopyInto(out *PodSpecData) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodSpecData.
func (in *PodSpecData) DeepCopy() *PodSpecData {
	if in == nil {
		return nil
	}
	out := new(PodSpecData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Run) DeepCopyInto(out *Run) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Run.
func (in *Run) DeepCopy() *Run {
	if in == nil {
		return nil
	}
	out := new(Run)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Run) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunList) DeepCopyInto(out *RunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Run, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunList.
func (in *RunList) DeepCopy() *RunList {
	if in == nil {
		return nil
	}
	out := new(RunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunSpec) DeepCopyInto(out *RunSpec) {
	*out = *in
	in.Build.DeepCopyInto(&out.Build)
	out.OutputData = in.OutputData
	out.InputData = in.InputData
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunSpec.
func (in *RunSpec) DeepCopy() *RunSpec {
	if in == nil {
		return nil
	}
	out := new(RunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RunStatus) DeepCopyInto(out *RunStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RunStatus.
func (in *RunStatus) DeepCopy() *RunStatus {
	if in == nil {
		return nil
	}
	out := new(RunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SourceCodeSpec) DeepCopyInto(out *SourceCodeSpec) {
	*out = *in
	in.Dependencies.DeepCopyInto(&out.Dependencies)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SourceCodeSpec.
func (in *SourceCodeSpec) DeepCopy() *SourceCodeSpec {
	if in == nil {
		return nil
	}
	out := new(SourceCodeSpec)
	in.DeepCopyInto(out)
	return out
}
