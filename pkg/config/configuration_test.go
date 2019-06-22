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
repository: myrepo.com
name: stuff
idl_directory: ./idl

dependencies:
  - name: dependency1
    version: 0.7.0
    type: protobuf

  - name: dependency2
    version: 2.9.4
    type: avro

provides:
  - root: ./docs/protos
    type: protobuf

  - root: ./docs/avros
    type: avro
`

		configuration := new(config.Configuration)

		err := configuration.UnMarshal(strings.NewReader(text))

		It("should not error", func() {
			Expect(err).To(BeNil())
		})

		It("should create settings object", func() {
			Expect(configuration.Repository).To(Equal("myrepo.com"))
			Expect(configuration.Name).To(Equal("stuff"))
			Expect(configuration.IdlDirectory).To(Equal("./idl"))
		})

		It("should create dependencies", func() {
			Expect(configuration.Dependencies).To(HaveLen(2))

			Expect(configuration.Dependencies[0]).
				To(Equal(config.Dependency{
					Name:    "dependency1",
					Version: "0.7.0",
					Type:    "protobuf",
				}))

			Expect(configuration.Dependencies[1]).
				To(Equal(config.Dependency{
					Name:    "dependency2",
					Version: "2.9.4",
					Type:    "avro",
				}))
		})

		It("should create provides", func() {
			Expect(configuration.Provides).To(HaveLen(2))

			Expect(configuration.Provides[0]).
				To(Equal(config.Provide{
					Root: "./docs/protos",
					Type: "protobuf",
				}))

			Expect(configuration.Provides[1]).
				To(Equal(config.Provide{
					Root: "./docs/avros",
					Type: "avro",
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
			Repository:   "our-company-repo.com",
			IdlDirectory: "./idl",
			Dependencies: []config.Dependency{
				config.Dependency{
					Name:    "dependency1",
					Version: "0.8.6",
					Type:    "protobuf",
				},
				config.Dependency{
					Name:    "dependency2",
					Version: "3.8.9",
					Type:    "avro",
				},
			},
			Provides: []config.Provide{
				config.Provide{
					Root: "./docs/proto",
					Type: "protobuf",
				},
				config.Provide{
					Root: "./docs/avro",
					Type: "avro",
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
			expect := `repository: our-company-repo.com
name: great-project
idl_directory: ./idl
dependencies:
- name: dependency1
  version: 0.8.6
  type: protobuf
- name: dependency2
  version: 3.8.9
  type: avro
provides:
- root: ./docs/proto
  type: protobuf
- root: ./docs/avro
  type: avro
`

			Expect(actual).To(Equal(expect))
		})
	})
})
