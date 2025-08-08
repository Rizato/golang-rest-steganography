package handlers

// These are request definitions, which should mostly match the model, just without the user specified fields

// TODO Algorithm choices

type CreateEmbedRequest struct {
	imageUUID string
	message   string
}

type CreateExtractRequest struct {
	imageUUID string
}
