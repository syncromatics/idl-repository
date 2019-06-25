package config_test

import (
	"bytes"
	"strings"

	"github.com/syncromatics/idl-repository/pkg/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	Context("unmarshalling from text", func() {
		text := `
repository: first.example.com
name: stuff
idl_directory: ./idl

dependencies:
  - name: dependency1
    version: 0.7.0
    type: protobuf

  - name: dependency2
    version: 2.9.4
    type: avro
    repository: second.example.net

provides:
  - root: ./docs/protos
    type: protobuf
    idlignore: custom_ignore_file

  - root: ./docs/avros
    type: avro
    idlignore: |-
        .noise
        .tmp
        *~
`

		configuration := new(config.Configuration)

		err := configuration.UnMarshal(strings.NewReader(text))

		It("should not error", func() {
			Expect(err).To(BeNil())
		})

		It("should create settings object", func() {
			Expect(configuration.Repository).To(Equal("first.example.com"))
			Expect(configuration.Name).To(Equal("stuff"))
			Expect(configuration.IdlDirectory).To(Equal("./idl"))
		})

		It("should create dependencies", func() {
			Expect(configuration.Dependencies).To(HaveLen(2))

			Expect(configuration.Dependencies[0]).
				To(Equal(config.Dependency{
					Name:       "dependency1",
					Version:    "0.7.0",
					Type:       "protobuf",
					Repository: "",
				}))

			Expect(configuration.Dependencies[1]).
				To(Equal(config.Dependency{
					Name:       "dependency2",
					Version:    "2.9.4",
					Type:       "avro",
					Repository: "second.example.net",
				}))
		})

		It("should create provides", func() {
			Expect(configuration.Provides).To(HaveLen(2))

			Expect(configuration.Provides[0]).
				To(Equal(config.Provide{
					Root:      "./docs/protos",
					Type:      "protobuf",
					IdlIgnore: "custom_ignore_file",
				}))

			Expect(configuration.Provides[1]).
				To(Equal(config.Provide{
					Root:      "./docs/avros",
					Type:      "avro",
					IdlIgnore: ".noise\n.tmp\n*~",
				}))
		})
	})

	Context("unmarshalling from bad text", func() {
		text := `this is clearly wrong`

		configuration := new(config.Configuration)

		err := configuration.UnMarshal(strings.NewReader(text))

		It("should error", func() {
			Expect(err).ToNot(BeNil())
		})
	})

	Context("marshalling configuration to text", func() {
		configuration := &config.Configuration{
			Name:         "great-project",
			Repository:   "first.example.com",
			IdlDirectory: "./idl",
			Dependencies: []config.Dependency{
				config.Dependency{
					Name:    "dependency1",
					Version: "0.8.6",
					Type:    "protobuf",
				},
				config.Dependency{
					Name:       "dependency2",
					Version:    "3.8.9",
					Type:       "avro",
					Repository: "second.example.net",
				},
			},
			Provides: []config.Provide{
				config.Provide{
					Root:      "./docs/proto",
					Type:      "protobuf",
					IdlIgnore: "custom_ignore_file",
				},
				config.Provide{
					Root:      "./docs/avro",
					Type:      "avro",
					IdlIgnore: ".noise\n.tmp\n*~",
				},
			},
		}

		text := new(bytes.Buffer)

		err := configuration.Marshal(text)

		It("should not error", func() {
			Expect(err).To(BeNil())
		})

		It("should create readable text", func() {
			actual := string(text.Bytes())
			expect := `repository: first.example.com
name: great-project
idl_directory: ./idl
dependencies:
- name: dependency1
  version: 0.8.6
  type: protobuf
  repository: ""
- name: dependency2
  version: 3.8.9
  type: avro
  repository: second.example.net
provides:
- root: ./docs/proto
  type: protobuf
  idlignore: custom_ignore_file
- root: ./docs/avro
  type: avro
  idlignore: |-
    .noise
    .tmp
    *~
`

			Expect(actual).To(Equal(expect))
		})
	})
})
