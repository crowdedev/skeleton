package configs

type (
	Model interface {
		TableName() string
		Identifier() string
		SetIdentifier(id string)
		SetCreatedBy(user *User)
		SetUpdatedBy(user *User)
		SetDeletedBy(user *User)
		IsSoftDelete() bool
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
