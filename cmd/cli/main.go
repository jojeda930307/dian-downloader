package main

import (
	"dian-downloader/internal/client"
	"flag" // Paquete estándar para flags de consola
	"fmt"
	"log"
	"os"
)

func main() {
	// 1. Definimos los parámetros (flags)
	// Sintaxis: flag.String("nombre", "valor_por_defecto", "descripción")
	docKey := flag.String("key", "", "El DocumentKey de la DIAN a descargar (obligatorio)")
	outputPath := flag.String("out", "", "Ruta y nombre del archivo PDF de salida (opcional)")

	// 2. Procesamos los argumentos de la consola
	flag.Parse()

	// 3. Validación de parámetros obligatorios
	if *docKey == "" {
		fmt.Println("Uso: ./dian-cli -key=\"TU_DOCUMENT_KEY\" [-out=\"ruta/archivo.pdf\"]")
		flag.PrintDefaults() // Imprime la ayuda automáticamente
		os.Exit(1)
	}

	// 4. Lógica de nombre de archivo por defecto si no se pasa -out
	finalPath := *outputPath
	if finalPath == "" {
		_ = os.MkdirAll("downloads", os.ModePerm) // Aseguramos que la carpeta exista
		finalPath = fmt.Sprintf("downloads/%s.pdf", *docKey)
	}

	log.Printf("[CLI] Iniciando descarga para: %s", *docKey)

	// 5. Instanciamos el cliente interno
	dian := client.NewDianClient()
	err := dian.DownloadPDF(*docKey, finalPath)

	if err != nil {
		log.Fatalf("[CLI ERROR] Fallo en la descarga: %v", err)
	}

	log.Printf("[CLI SUCCESS] Proceso finalizado. PDF guardado en: %s", finalPath)
}
