package internal

type CreateMicroserviceRequest struct {
	Name        string `json:"name" binding:"required" example:"api-flask"`
	Description string `json:"description" example:"API de prueba en Python"`
	Language    string `json:"language" binding:"required" example:"flask"`
	Code        string `json:"code" binding:"required" example:"from flask import Flask\napp = Flask(__name__)\n\n@app.route('/')\ndef hello():\n    return 'Hola Mundo!'\n\nif __name__ == '__main__':\n    app.run(host='0.0.0.0')"`
}

type UpdateMicroserviceRequest struct {
	Code        *string `json:"code" example:"from flask import Flask\napp = Flask(__name__)\n\n@app.route('/')\ndef hello():\n    return 'Adios Mundo!'\n\nif __name__ == '__main__':\n    app.run(host='0.0.0.0')"`
	Description *string `json:"description" example:"API actualizada de prueba"`
}
