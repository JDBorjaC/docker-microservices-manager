from flask import Flask, jsonify, request

app = Flask(__name__)

@app.route('/saludo', methods=['POST'])
def saludar():
    datos = request.get_json()
    nombre = datos.get('nombre', 'Invitado')
    return jsonify({'mensaje': f'Hola {nombre}!', 'status': 'success'})

if __name__ == '__main__':
    app.run(debug=True)