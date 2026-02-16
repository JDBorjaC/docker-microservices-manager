import { useState, useEffect } from "react";
import "./styles/index.css";

// Mock data
const mockPreguntas = {
    numPreg: 5,
    enunciados: [
        { enunciado: "¿Cuál es el nombre del usuario?" },
        { enunciado: "¿Cuál es la edad del usuario?" },
        { enunciado: "Escribe una consulta SQL para obtener usuarios mayores de edad" },
        { enunciado: "¿Cuál es la edad del usuario?" },
        { enunciado: "Escribe una consulta SQL para obtener usuarios mayores de edad" }
    ]
};

const mockTabla = [
    { nombre: "Pepito Perez", edad: 2 },
    { nombre: "Juan Lopez", edad: 20 }
];

export default function Examen() {
    const [preguntaActiva, setPreguntaActiva] = useState(0);
    const [respuestas, setRespuestas] = useState({});
    const enunciados = mockPreguntas.enunciados;

    // --- Cargar respuestas guardadas ---
    useEffect(() => {
        const stored = localStorage.getItem("respuestasExamen");
        if (stored) setRespuestas(JSON.parse(stored));
    }, []);

    // --- Guardar respuesta ---
    const handleRespuestaChange = (texto) => {
        const nuevas = {
            ...respuestas,
            [preguntaActiva]: texto,
        };

        setRespuestas(nuevas);
        localStorage.setItem("respuestasExamen", JSON.stringify(nuevas));
    };

    const respuestaActual = respuestas[preguntaActiva] || "";

    // --- Enviar examen y limpiar todo ---
    const enviarExamen = () => {
        alert("Examen enviado:\n" + JSON.stringify(respuestas, null, 2));

        // Limpiar localStorage
        localStorage.removeItem("respuestasExamen");

        // Limpiar estado
        setRespuestas({});
    };

    return (
        <div className="examen-container">

            {/* Panel de preguntas */}
            <div className="preguntas-panel">
                <div className="preguntas-tabs">
                    {enunciados.map((_, index) => (
                        <button
                            key={index}
                            className={
                                "pregunta-tab " + (preguntaActiva === index ? "active" : "")
                            }
                            onClick={() => setPreguntaActiva(index)}
                        >
                            Pregunta {index + 1}
                        </button>
                    ))}
                </div>

                <div className="pregunta-enunciado">
                    {enunciados[preguntaActiva].enunciado}
                </div>

                <div className="pregunta-tabla">
                    <table>
                        <thead>
                            <tr>
                                <th>Nombre</th>
                                <th>Edad</th>
                            </tr>
                        </thead>
                        <tbody>
                            {mockTabla.map((fila, index) => (
                                <tr key={index}>
                                    <td>{fila.nombre}</td>
                                    <td>{fila.edad}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Panel SQL */}
            <div className="sql-panel">
                <textarea
                    className="sql-input"
                    placeholder="Escribe tu código SQL aquí"
                    value={respuestaActual}
                    onChange={(e) => handleRespuestaChange(e.target.value)}
                />

                <div className="buttons">
                    <button className="btn-primary" onClick={enviarExamen}>
                        Enviar Examen
                    </button>
                    <button className="btn-secondary" onClick>
                        Probar Consulta
                    </button>
                </div>
            </div>
        </div>
    );
}
