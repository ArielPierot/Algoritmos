package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type pacoteStruct struct {
	ordem     int
	tamPacote int
	pacotes   []string
}

func main() {
	aquivoEntrada := os.Args[1] // Entrada
	arquivoSaida := os.Args[2]  // Saída

	// Leitura de arquivo
	arquivo, err := os.Open(aquivoEntrada)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de entrada")
		return
	}
	scanner := bufio.NewScanner(arquivo)
	f, errOut := os.Create(arquivoSaida)

	var arrayPacotes []*pacoteStruct

	if scanner.Scan() {

		linha := scanner.Text()
		palavra := strings.Fields(linha)

		numeroTotalPacotes, _ := strconv.Atoi(palavra[0])
		quantidadePacotes, _ := strconv.Atoi(palavra[1])

		for i := 0; i < numeroTotalPacotes && scanner.Scan(); i++ {

			var arr []string

			linha := scanner.Text()
			palavra := strings.Fields(linha)

			ordem, _ := strconv.Atoi(palavra[0])
			tamanhoPacotes, _ := strconv.Atoi(palavra[1])

			for n := 2; n <= tamanhoPacotes+1; n++ {
				arr = append(arr, palavra[n])
			}

			arrayPacotes = append(arrayPacotes, &pacoteStruct{ordem, tamanhoPacotes, arr})
		}

		// Lógica da transmissão
		if errOut == nil {
			var numberConsecutivo, inicioTransmissao, lastConsecutivo, i int
			var transmitir bool

			w := bufio.NewWriter(f)

			var divisaoFloat float64
			var tempResto float64
			divisaoFloat = float64(numeroTotalPacotes) / float64(quantidadePacotes)
			divisaoInt := numeroTotalPacotes / quantidadePacotes
			if divisaoFloat > float64(divisaoInt) {
				tempResto = divisaoFloat - float64(divisaoInt)
				tempResto = tempResto * float64(quantidadePacotes)
			}

			for i = 0; i < (numeroTotalPacotes / quantidadePacotes); i++ {

				heapsort(arrayPacotes, quantidadePacotes*(i+1))

				transmitir = existsConsutivo(arrayPacotes, quantidadePacotes*(i+1), numberConsecutivo)

				if transmitir {

					lastConsecutivo = numberConsecutivo
					numberConsecutivo = 0

					transmissao := i - inicioTransmissao
					w.WriteString(fmt.Sprintf("%d:", transmissao))

					for l := 0; l < (quantidadePacotes * (i + 1)); l++ {

						if arrayPacotes[l].ordem == numberConsecutivo {
							if arrayPacotes[l].ordem >= lastConsecutivo {
								for _, v := range arrayPacotes[l].pacotes {
									w.WriteString(fmt.Sprintf(" %s", v))
								}
							}
						} else {
							break
						}
						numberConsecutivo = arrayPacotes[l].ordem + 1

					}

					w.WriteString(fmt.Sprintf("\r\n"))

				} else {
					inicioTransmissao++
				}

			}

			// Se houver resto, então ele transmite todos os pacotes ordenados a partir da última posição de ordem
			if tempResto > 0 {
				heapsort(arrayPacotes, numeroTotalPacotes)
				transmissao := i - inicioTransmissao
				w.WriteString(fmt.Sprintf("%d:", transmissao))

				for l := 0; l < numeroTotalPacotes; l++ {
					if arrayPacotes[l].ordem >= numberConsecutivo {
						for _, v := range arrayPacotes[l].pacotes {
							w.WriteString(fmt.Sprintf(" %s", v))
						}
					}
				}

				w.WriteString(fmt.Sprintf("\r\n"))
			}

			w.Flush()
			arquivo.Close()
			f.Close()
		}
	}
}

func esquerdo(pai int) int {
	return 2*pai + 1
}

func direito(pai int) int {
	return 2*pai + 2
}

func heapify(pacote []*pacoteStruct, i int, n int) {
	P := i
	E := esquerdo(i)
	D := direito(i)
	if E < n && pacote[E].ordem > pacote[P].ordem {
		P = E
	}

	if D < n && pacote[D].ordem > pacote[P].ordem {
		P = D
	}

	if P != i {
		pacote[P], pacote[i] = pacote[i], pacote[P]
		heapify(pacote, P, n)
	}
}

func maxheap(pacote []*pacoteStruct, n int) {
	for i := n / 2; i >= 0; i-- {
		heapify(pacote, i, n)
	}
}

func heapsort(pacote []*pacoteStruct, n int) {
	maxheap(pacote, n)
	var temp *pacoteStruct

	for i := n - 1; i > 0; i-- {
		temp = pacote[0]
		pacote[0] = pacote[i]
		pacote[i] = temp
		heapify(pacote, 0, i)
	}
}

func existsConsutivo(pacote []*pacoteStruct, n int, consecutivo int) bool {
	for i := 0; i < n; i++ {
		if pacote[i].ordem == consecutivo {
			return true
		}
	}
	return false
}
