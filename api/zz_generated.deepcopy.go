//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package api

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *App) DeepCopyInto(out *App) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new App.
func (in *App) DeepCopy() *App {
	if in == nil {
		return nil
	}
	out := new(App)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppDetails) DeepCopyInto(out *AppDetails) {
	*out = *in
	out.App = in.App
	out.Details = in.Details
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppDetails.
func (in *AppDetails) DeepCopy() *AppDetails {
	if in == nil {
		return nil
	}
	out := new(AppDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentDefinition) DeepCopyInto(out *ComponentDefinition) {
	*out = *in
	if in.Routes != nil {
		in, out := &in.Routes, &out.Routes
		*out = make([]RouteSpec, len(*in))
		copy(*out, *in)
	}
	if in.EnvSchema != nil {
		in, out := &in.EnvSchema, &out.EnvSchema
		*out = make(map[string]*EnvVarSchema, len(*in))
		for key, val := range *in {
			var outVal *EnvVarSchema
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(EnvVarSchema)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.Dependencies != nil {
		in, out := &in.Dependencies, &out.Dependencies
		*out = make(map[string]*Dependency, len(*in))
		for key, val := range *in {
			var outVal *Dependency
			if val == nil {
				(*out)[key] = nil
			} else {
				inVal := (*in)[key]
				in, out := &inVal, &outVal
				*out = new(Dependency)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentDefinition.
func (in *ComponentDefinition) DeepCopy() *ComponentDefinition {
	if in == nil {
		return nil
	}
	out := new(ComponentDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentDetails) DeepCopyInto(out *ComponentDetails) {
	*out = *in
	in.ComponentDefinition.DeepCopyInto(&out.ComponentDefinition)
	out.Details = in.Details
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentDetails.
func (in *ComponentDetails) DeepCopy() *ComponentDetails {
	if in == nil {
		return nil
	}
	out := new(ComponentDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dependency) DeepCopyInto(out *Dependency) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dependency.
func (in *Dependency) DeepCopy() *Dependency {
	if in == nil {
		return nil
	}
	out := new(Dependency)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Details) DeepCopyInto(out *Details) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Details.
func (in *Details) DeepCopy() *Details {
	if in == nil {
		return nil
	}
	out := new(Details)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvVarSchema) DeepCopyInto(out *EnvVarSchema) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvVarSchema.
func (in *EnvVarSchema) DeepCopy() *EnvVarSchema {
	if in == nil {
		return nil
	}
	out := new(EnvVarSchema)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RouteSpec) DeepCopyInto(out *RouteSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RouteSpec.
func (in *RouteSpec) DeepCopy() *RouteSpec {
	if in == nil {
		return nil
	}
	out := new(RouteSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Val) DeepCopyInto(out *Val) {
	*out = *in
	if in.arrayNumVal != nil {
		in, out := &in.arrayNumVal, &out.arrayNumVal
		*out = make([]float64, len(*in))
		copy(*out, *in)
	}
	if in.arrayStrVal != nil {
		in, out := &in.arrayStrVal, &out.arrayStrVal
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Val.
func (in *Val) DeepCopy() *Val {
	if in == nil {
		return nil
	}
	out := new(Val)
	in.DeepCopyInto(out)
	return out
}