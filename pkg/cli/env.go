package cli

import (
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	namespace string
	config    *genericclioptions.ConfigFlags

	// KubeConfig is the path to the kubeconfig file
	KubeConfig string
	// KubeContext is the name of the kubeconfig context.
	KubeContext string
	// Bearer KubeToken used for authentication
	LogLevel string
}

func New() *EnvSettings {
	env := &EnvSettings{
		namespace:   os.Getenv("TB_NAMESPACE"),
		KubeContext: os.Getenv("TB_KUBECONTEXT"),
		KubeConfig:  os.Getenv("TB_KUBECONFIG"),
		LogLevel:    envOr("TB_LOGLEVEL", "info"),
	}

	// bind to kubernetes config flags
	env.config = &genericclioptions.ConfigFlags{
		Namespace:  &env.namespace,
		Context:    &env.KubeContext,
		KubeConfig: &env.KubeConfig,
	}
	return env
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&s.namespace, "namespace", "n", s.namespace, "namespace scope for this request")
	fs.StringVar(&s.KubeConfig, "kubeconfig", s.KubeConfig, "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVarP(&s.LogLevel, "loglevel", "l", s.LogLevel, "Log level")
}

func envOr(name, def string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return def
}

func (s *EnvSettings) EnvVars() map[string]string {
	envvars := map[string]string{
		"TB_BIN":         os.Args[0],
		"TB_LOGLEVEL":    s.LogLevel,
		"TB_NAMESPACE":   s.Namespace(),
		"TB_KUBECONTEXT": s.KubeContext,
		"TB_KUBECONFIG":  s.KubeConfig,
	}
	if s.KubeConfig != "" {
		envvars["KUBECONFIG"] = s.KubeConfig
	}
	return envvars
}

// Namespace gets the namespace from the configuration
func (s *EnvSettings) Namespace() string {
	if ns, _, err := s.config.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

// SetNamespace sets the namespace in the configuration
func (s *EnvSettings) SetNamespace(namespace string) {
	s.namespace = namespace
}

// RESTClientGetter gets the kubeconfig from EnvSettings
func (s *EnvSettings) RESTClientGetter() genericclioptions.RESTClientGetter {
	return s.config
}
