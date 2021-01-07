package configs

type (
	Model interface {
		TableName() string
		Identifier() string
		SetIdentifier(id string)
	}

	Service interface {
		Model() Model
		Create() Model
		Update() Model
		Bind() Model
		Delete()
	}

	Application interface {
		Run()
	}
)
