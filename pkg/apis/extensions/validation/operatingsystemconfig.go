// SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateOperatingSystemConfig validates a OperatingSystemConfig object.
func ValidateOperatingSystemConfig(osc *extensionsv1alpha1.OperatingSystemConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, apivalidation.ValidateObjectMeta(&osc.ObjectMeta, true, apivalidation.NameIsDNSSubdomain, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateOperatingSystemConfigSpec(&osc.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateOperatingSystemConfigUpdate validates a OperatingSystemConfig object before an update.
func ValidateOperatingSystemConfigUpdate(new, old *extensionsv1alpha1.OperatingSystemConfig) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateObjectMetaUpdate(&new.ObjectMeta, &old.ObjectMeta, field.NewPath("metadata"))...)
	allErrs = append(allErrs, ValidateOperatingSystemConfigSpecUpdate(&new.Spec, &old.Spec, new.DeletionTimestamp != nil, field.NewPath("spec"))...)
	allErrs = append(allErrs, ValidateOperatingSystemConfig(new)...)

	return allErrs
}

// ValidateOperatingSystemConfigSpec validates the specification of a OperatingSystemConfig object.
func ValidateOperatingSystemConfigSpec(spec *extensionsv1alpha1.OperatingSystemConfigSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(spec.Type) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("type"), "field is required"))
	}

	if len(spec.Purpose) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("purpose"), "field is required"))
	} else {
		if spec.Purpose != extensionsv1alpha1.OperatingSystemConfigPurposeProvision && spec.Purpose != extensionsv1alpha1.OperatingSystemConfigPurposeReconcile {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("purpose"), spec.Purpose, []string{string(extensionsv1alpha1.OperatingSystemConfigPurposeProvision), string(extensionsv1alpha1.OperatingSystemConfigPurposeReconcile)}))
		}
	}

	allErrs = append(allErrs, ValidateUnits(spec.Units, fldPath.Child("units"))...)
	allErrs = append(allErrs, ValidateFiles(spec.Files, fldPath.Child("files"))...)

	return allErrs
}

// ValidateUnits validates operating system config units.
func ValidateUnits(units []extensionsv1alpha1.Unit, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, unit := range units {
		idxPath := fldPath.Index(i)

		if len(unit.Name) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), "field is required"))
		}

		for j, dropIn := range unit.DropIns {
			jdxPath := idxPath.Child("dropIns").Index(j)

			if len(dropIn.Name) == 0 {
				allErrs = append(allErrs, field.Required(jdxPath.Child("name"), "field is required"))
			}
			if len(dropIn.Content) == 0 {
				allErrs = append(allErrs, field.Required(jdxPath.Child("content"), "field is required"))
			}
		}
	}

	return allErrs
}

// ValidateFiles validates operating system config files.
func ValidateFiles(files []extensionsv1alpha1.File, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	paths := sets.NewString()

	for i, file := range files {
		idxPath := fldPath.Index(i)

		if len(file.Path) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("path"), "field is required"))
		} else {
			if paths.Has(file.Path) {
				allErrs = append(allErrs, field.Duplicate(idxPath.Child("path"), file.Path))
			}
			paths.Insert(file.Path)
		}

		switch {
		case file.Content.SecretRef == nil && file.Content.Inline == nil:
			allErrs = append(allErrs, field.Required(idxPath.Child("content"), "either 'secretRef' or 'inline' must be provided"))
		case file.Content.SecretRef != nil && file.Content.Inline != nil:
			allErrs = append(allErrs, field.Invalid(idxPath.Child("content"), file.Content, "either 'secretRef' or 'inline' must be provided, not both at the same time"))
		case file.Content.SecretRef != nil:
			if len(file.Content.SecretRef.Name) == 0 {
				allErrs = append(allErrs, field.Required(idxPath.Child("content", "secretRef", "name"), "field is required"))
			}
			if len(file.Content.SecretRef.DataKey) == 0 {
				allErrs = append(allErrs, field.Required(idxPath.Child("content", "secretRef", "dataKey"), "field is required"))
			}
		case file.Content.Inline != nil:
			if len(file.Content.Inline.Encoding) == 0 {
				allErrs = append(allErrs, field.Required(idxPath.Child("content", "inline", "encoding"), "field is required"))
			}
			if len(file.Content.Inline.Data) == 0 {
				allErrs = append(allErrs, field.Required(idxPath.Child("content", "inline", "data"), "field is required"))
			}
		}
	}

	return allErrs
}

// ValidateOperatingSystemConfigSpecUpdate validates the spec of a OperatingSystemConfig object before an update.
func ValidateOperatingSystemConfigSpecUpdate(new, old *extensionsv1alpha1.OperatingSystemConfigSpec, deletionTimestampSet bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if deletionTimestampSet && !apiequality.Semantic.DeepEqual(new, old) {
		allErrs = append(allErrs, apivalidation.ValidateImmutableField(new, old, fldPath)...)
		return allErrs
	}

	allErrs = append(allErrs, apivalidation.ValidateImmutableField(new.Type, old.Type, fldPath.Child("type"))...)
	allErrs = append(allErrs, apivalidation.ValidateImmutableField(new.Purpose, old.Purpose, fldPath.Child("purpose"))...)

	return allErrs
}

// ValidateOperatingSystemConfigStatus validates the status of a OperatingSystemConfig object.
func ValidateOperatingSystemConfigStatus(spec *extensionsv1alpha1.OperatingSystemConfigStatus, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}

// ValidateOperatingSystemConfigStatusUpdate validates the status field of a OperatingSystemConfig object.
func ValidateOperatingSystemConfigStatusUpdate(newStatus, oldStatus extensionsv1alpha1.OperatingSystemConfigStatus) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}
