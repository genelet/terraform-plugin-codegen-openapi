// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oas_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-codegen-openapi/internal/mapper/oas"
	"github.com/hashicorp/terraform-plugin-codegen-spec/datasource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/resource"
	"github.com/hashicorp/terraform-plugin-codegen-spec/schema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// TODO: add error tests

func TestBuildMapResource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]resource.Attribute
	}{
		"map nested attribute with props": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map nested type.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type:     []string{"object"},
							Required: []string{"nested_password_required"},
							Properties: map[string]*base.SchemaProxy{
								"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"object"},
									Required:    []string{"nested_int64_required"},
									Description: "hey there! I'm a single nested object type.",
									Properties: map[string]*base.SchemaProxy{
										"nested_float64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"number"},
											Format:      "double",
											Description: "hey there! I'm a nested float64 type.",
										}),
										"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"integer"},
											Format:      "int64",
											Description: "hey there! I'm a nested int64 type, required.",
										}),
									},
								}),
								"nested_password_required": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"string"},
									Format:      "password",
									Description: "hey there! I'm a nested string type, required.",
								}),
							},
						}),
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "nested_map_prop",
					MapNested: &resource.MapNestedAttribute{
						NestedObject: resource.NestedAttributeObject{
							Attributes: []resource.Attribute{
								{
									Name: "nested_obj_prop",
									SingleNested: &resource.SingleNestedAttribute{
										Attributes: []resource.Attribute{
											{
												Name: "nested_float64",
												Float64: &resource.Float64Attribute{
													ComputedOptionalRequired: schema.ComputedOptional,
													Description:              pointer("hey there! I'm a nested float64 type."),
												},
											},
											{
												Name: "nested_int64_required",
												Int64: &resource.Int64Attribute{
													ComputedOptionalRequired: schema.Required,
													Description:              pointer("hey there! I'm a nested int64 type, required."),
												},
											},
										},
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a single nested object type."),
									},
								},
								{
									Name: "nested_password_required",
									String: &resource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Sensitive:                pointer(true),
										Description:              pointer("hey there! I'm a nested string type, required."),
									},
								},
							},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map nested type."),
					},
				},
			},
		},
		"map attributes with element types": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"map_with_strings_required"},
				Properties: map[string]*base.SchemaProxy{
					"map_with_floats": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with floats.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type:   []string{"number"},
							Format: "float",
						}),
					}),
					"map_with_strings_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with strings, required.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type: []string{"string"},
						}),
					}),
				},
			},
			expectedAttributes: &[]resource.Attribute{
				{
					Name: "map_with_floats",
					Map: &resource.MapAttribute{
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map type with floats."),
					},
				},
				{
					Name: "map_with_strings_required",
					Map: &resource.MapAttribute{
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a map type with strings, required."),
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := oas.OASSchema{Schema: testCase.schema}
			attributes, err := schema.BuildResourceAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestBuildMapDataSource(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		schema             *base.Schema
		expectedAttributes *[]datasource.Attribute
	}{

		"map nested attribute with props": {
			schema: &base.Schema{
				Type: []string{"object"},
				Properties: map[string]*base.SchemaProxy{
					"nested_map_prop": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map nested type.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type:     []string{"object"},
							Required: []string{"nested_password_required"},
							Properties: map[string]*base.SchemaProxy{
								"nested_obj_prop": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"object"},
									Required:    []string{"nested_int64_required"},
									Description: "hey there! I'm a single nested object type.",
									Properties: map[string]*base.SchemaProxy{
										"nested_float64": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"number"},
											Format:      "double",
											Description: "hey there! I'm a nested float64 type.",
										}),
										"nested_int64_required": base.CreateSchemaProxy(&base.Schema{
											Type:        []string{"integer"},
											Format:      "int64",
											Description: "hey there! I'm a nested int64 type, required.",
										}),
									},
								}),
								"nested_password_required": base.CreateSchemaProxy(&base.Schema{
									Type:        []string{"string"},
									Format:      "password",
									Description: "hey there! I'm a nested string type, required.",
								}),
							},
						}),
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "nested_map_prop",
					MapNested: &datasource.MapNestedAttribute{
						NestedObject: datasource.NestedAttributeObject{
							Attributes: []datasource.Attribute{
								{
									Name: "nested_obj_prop",
									SingleNested: &datasource.SingleNestedAttribute{
										Attributes: []datasource.Attribute{
											{
												Name: "nested_float64",
												Float64: &datasource.Float64Attribute{
													ComputedOptionalRequired: schema.ComputedOptional,
													Description:              pointer("hey there! I'm a nested float64 type."),
												},
											},
											{
												Name: "nested_int64_required",
												Int64: &datasource.Int64Attribute{
													ComputedOptionalRequired: schema.Required,
													Description:              pointer("hey there! I'm a nested int64 type, required."),
												},
											},
										},
										ComputedOptionalRequired: schema.ComputedOptional,
										Description:              pointer("hey there! I'm a single nested object type."),
									},
								},
								{
									Name: "nested_password_required",
									String: &datasource.StringAttribute{
										ComputedOptionalRequired: schema.Required,
										Sensitive:                pointer(true),
										Description:              pointer("hey there! I'm a nested string type, required."),
									},
								},
							},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map nested type."),
					},
				},
			},
		},
		"map attributes with element types": {
			schema: &base.Schema{
				Type:     []string{"object"},
				Required: []string{"map_with_strings_required"},
				Properties: map[string]*base.SchemaProxy{
					"map_with_floats": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with floats.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type:   []string{"number"},
							Format: "float",
						}),
					}),
					"map_with_strings_required": base.CreateSchemaProxy(&base.Schema{
						Type:        []string{"object"},
						Description: "hey there! I'm a map type with strings, required.",
						AdditionalProperties: base.CreateSchemaProxy(&base.Schema{
							Type: []string{"string"},
						}),
					}),
				},
			},
			expectedAttributes: &[]datasource.Attribute{
				{
					Name: "map_with_floats",
					Map: &datasource.MapAttribute{
						ElementType: schema.ElementType{
							Float64: &schema.Float64Type{},
						},
						ComputedOptionalRequired: schema.ComputedOptional,
						Description:              pointer("hey there! I'm a map type with floats."),
					},
				},
				{
					Name: "map_with_strings_required",
					Map: &datasource.MapAttribute{
						ElementType: schema.ElementType{
							String: &schema.StringType{},
						},
						ComputedOptionalRequired: schema.Required,
						Description:              pointer("hey there! I'm a map type with strings, required."),
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			schema := oas.OASSchema{Schema: testCase.schema}
			attributes, err := schema.BuildDataSourceAttributes()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if diff := cmp.Diff(attributes, testCase.expectedAttributes); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
