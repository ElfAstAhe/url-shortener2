package config

// DefaultConfFactory - фабрика экземпляра default конфигурации
type DefaultConfFactory func() *AppConf

type ParamsInitializer func([]*AppConfSource) error

type Loader func(params ...any) error

// AppConfSource - источники конфигурации
type AppConfSource struct {
	name           string
	priority       int
	defConfFactory DefaultConfFactory
	conf           *AppConf
	params         []any
	paramsInit     ParamsInitializer
	loader         Loader
}

func NewAppConfSource(
	name string,
	priority int,
	defConfFactory DefaultConfFactory,
	loader Loader,
) *AppConfSource {
	return NewAppConfSourceWithParams(name, priority, defConfFactory, nil, loader)
}

func NewAppConfSourceWithParams(
	name string,
	priority int,
	defConfFactory DefaultConfFactory,
	paramsInit ParamsInitializer,
	loader Loader,
) *AppConfSource {
	return &AppConfSource{
		name:           name,
		priority:       priority,
		defConfFactory: defConfFactory,
		params:         make([]any, 0),
		paramsInit:     paramsInit,
		loader:         loader,
	}
}

func (s *AppConfSource) Name() string {
	return s.name
}

func (s *AppConfSource) Priority() int {
	return s.priority
}

type AppConfLoader struct {
	confSources []*AppConfSource
}

func NewAppConfLoader(sources []*AppConfSource) *AppConfLoader {
	return &AppConfLoader{
		confSources: make([]*AppConfSource, 0),
	}
}

func (acl *AppConfLoader) LoadConfig() error {
	for _, source := range acl.confSources {
		if source.paramsInit != nil {
			err := source.paramsInit(acl.confSources)
			if err != nil {
				return err
			}
		}
		err := source.loader(source.params...)
		if err != nil {
			return err
		}
	}

	return nil
}
