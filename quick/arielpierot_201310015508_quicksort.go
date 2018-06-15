package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type vetorChamada struct {
	tipo  string
	calls int
}

func main() {

	arquivoEntrada := os.Args[1] // Entrada
	arquivoSaida := os.Args[2]   // Saída

	// Leitura de arquivo
	arquivo, err := os.Open(arquivoEntrada)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de entrada")
		return
	}
	scanner := bufio.NewScanner(arquivo)
	f, errOut := os.Create(arquivoSaida)

	if scanner.Scan() {
		totalVetores, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Erro ao ler linha => total de vetores")
		} else {
			for i := 0; i < totalVetores && scanner.Scan(); i++ {
				tamanhoVetor, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Println("Erro ao ler o tamanho do vetor => ", tamanhoVetor)
				} else {
					var vetorAux []int
					scanner.Scan()
					linha := scanner.Text()
					palavra := strings.Fields(linha)

					if len(palavra) != tamanhoVetor {
						fmt.Printf("Vetor está com tamanho incorreto. Valor esperado => %d. Valor verificado => %d \r\n", tamanhoVetor, len(palavra))
						break
					} else {
						for i, v := range palavra {
							vInt, err := strconv.Atoi(v)
							if err == nil {
								vetorAux = append(vetorAux, vInt)
							} else {
								fmt.Println("Erro ao fazer conversão de inteiro. Append do vetorAux => ", i)
							}
						}
						var vetorChamadas []*vetorChamada
						tamVetor := len(vetorAux) - 1

						tmp := make([]int, len(vetorAux))
						callsPP := 0
						copy(tmp, vetorAux)
						quicksort(tmp, 0, tamVetor, &callsPP)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"PP", callsPP})

						tmp = make([]int, len(vetorAux))
						copy(tmp, vetorAux)
						callsPM := 0
						quicksortMediana(tmp, 0, tamVetor, &callsPM)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"PM", callsPM})

						tmp = make([]int, len(vetorAux))
						copy(tmp, vetorAux)
						callsPA := 0
						quicksortRandom(tmp, 0, tamVetor, &callsPA)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"PA", callsPA})

						tmp = make([]int, len(vetorAux))
						copy(tmp, vetorAux)
						callsHP := 0
						hoare(tmp, 0, tamVetor, &callsHP)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"HP", callsHP})

						tmp = make([]int, len(vetorAux))
						copy(tmp, vetorAux)
						callsHM := 0
						hoareMediana(tmp, 0, tamVetor, &callsHM)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"HM", callsHM})

						tmp = make([]int, len(vetorAux))
						copy(tmp, vetorAux)
						callsHA := 0
						hoareRandom(tmp, 0, tamVetor, &callsHA)
						vetorChamadas = append(vetorChamadas, &vetorChamada{"HA", callsHA})

						vetorChamadasOrdenado := mergesort(vetorChamadas)

						if errOut == nil {
							w := bufio.NewWriter(f)
							w.WriteString(fmt.Sprintf("%d: N(%d)", i, len(vetorAux)))
							for _, v := range vetorChamadasOrdenado {
								w.WriteString(fmt.Sprintf(" %s(%d)", v.tipo, v.calls))
							}
							w.WriteString(fmt.Sprintf("\r\n"))
							w.Flush()
						} else {
							fmt.Println("Houve um erro na impressão dos dados")
						}

					}
				}
			}
			arquivo.Close()
		}
	}
}

func quicksort(vetor []int, inicio int, fim int, calls *int) {

	*calls++
	if inicio < fim {
		pivo := particionamento(vetor, inicio, fim, calls)
		quicksort(vetor, inicio, pivo-1, calls)
		quicksort(vetor, pivo+1, fim, calls)
	}
}

func quicksortRandom(vetor []int, inicio int, fim int, calls *int) {

	*calls++
	if inicio < fim {
		pivo_indice := randomPivo(vetor, inicio, fim, calls)
		quicksortRandom(vetor, inicio, pivo_indice-1, calls)
		quicksortRandom(vetor, pivo_indice+1, fim, calls)
	}
}

func quicksortMediana(vetor []int, inicio int, fim int, calls *int) {

	*calls++
	if inicio < fim {
		pivo_indice := medianaPivo(vetor, inicio, fim, calls)
		quicksortMediana(vetor, inicio, pivo_indice-1, calls)
		quicksortMediana(vetor, pivo_indice+1, fim, calls)
	}

}

func particionamento(vetor []int, inicio int, fim int, calls *int) int {

	pivo := vetor[fim]
	pivo_indice := inicio

	for i := inicio; i < fim; i++ {
		if vetor[i] <= pivo {
			vetor[i], vetor[pivo_indice] = vetor[pivo_indice], vetor[i]
			pivo_indice++
			*calls++
		}
	}

	vetor[pivo_indice], vetor[fim] = vetor[fim], vetor[pivo_indice]
	*calls++

	return pivo_indice
}

func medianaPivo(vetor []int, inicio int, fim int, calls *int) int {

	//*calls++

	n := len(vetor[inicio:fim]) + 1 // Corta o vetor do inicio ao fim e retorna o tamanho

	inicial := inicio + n/4
	meio := inicio + n/2
	final := inicio + 3*n/4

	var arr []int
	arr = append(arr, vetor[inicial])
	arr = append(arr, vetor[meio])
	arr = append(arr, vetor[final])

	var aux int
	quicksort(arr, 0, 2, &aux)

	if vetor[inicial] == arr[1] {
		vetor[inicial], vetor[fim] = vetor[fim], vetor[inicial]
	} else if vetor[meio] == arr[1] {
		vetor[meio], vetor[fim] = vetor[fim], vetor[meio]
	} else {
		vetor[final], vetor[fim] = vetor[fim], vetor[final]
	}

	*calls++

	return particionamento(vetor, inicio, fim, calls)
}

func randomPivo(vetor []int, inicio int, fim int, calls *int) int {

	n := len(vetor[inicio:fim]) + 1       // Corta o vetor do inicio ao fim e retorna o tamanho
	y := math.Abs(float64(vetor[inicio])) // Transforma em vetor[inicio] positivo
	resto := int(y) % n                   // Retorna o resto
	soma := (inicio + resto)              // Técnica para definição do índice, vai retornar um valor inteiro
	pivo_indice := soma                   // Escolha do pivô somando x da técnica a partir do ínicio

	vetor[pivo_indice], vetor[fim] = vetor[fim], vetor[pivo_indice] // Ocorre a troca do pivo pelo fim
	*calls++

	return particionamento(vetor, inicio, fim, calls) // Chama o método de particionar
}

func particionamentoHoare(vetor []int, inicio int, fim int, calls *int) int {

	pivo := vetor[inicio]
	i := inicio
	j := fim

	for i < j {
		for j > i && vetor[j] >= pivo {
			j--
		}
		for i < j && vetor[i] < pivo {
			i++
		}

		if i < j {
			vetor[i], vetor[j] = vetor[j], vetor[i]
			*calls++
		}
	}

	return j
}

func hoare(vetor []int, inicio int, fim int, calls *int) {
	*calls++

	if inicio < fim {
		pivo := particionamentoHoare(vetor, inicio, fim, calls)
		hoare(vetor, inicio, pivo, calls)
		hoare(vetor, pivo+1, fim, calls)
	}
}

func hoareMediana(vetor []int, inicio int, fim int, calls *int) {
	*calls++

	if inicio < fim {
		pivo := medianaPivoHoare(vetor, inicio, fim, calls)
		hoareMediana(vetor, inicio, pivo, calls)
		hoareMediana(vetor, pivo+1, fim, calls)
	}
}

func medianaPivoHoare(vetor []int, inicio int, fim int, calls *int) int {

	n := len(vetor[inicio:fim]) + 1 // Corta o vetor do inicio ao fim e retorna o tamanho

	inicial := inicio + n/4
	meio := inicio + n/2
	final := inicio + 3*n/4

	var arr []int
	arr = append(arr, vetor[inicial])
	arr = append(arr, vetor[meio])
	arr = append(arr, vetor[final])

	var aux int
	quicksort(arr, 0, 2, &aux)

	if vetor[inicial] == arr[1] {
		vetor[inicial], vetor[inicio] = vetor[inicio], vetor[inicial]
	} else if vetor[meio] == arr[1] {
		vetor[meio], vetor[inicio] = vetor[inicio], vetor[meio]
	} else {
		vetor[final], vetor[inicio] = vetor[inicio], vetor[final]
	}

	*calls++

	return particionamentoHoare(vetor, inicio, fim, calls)
}

func randomPivoHoare(vetor []int, inicio int, fim int, calls *int) int {

	n := len(vetor[inicio:fim]) + 1       // Corta o vetor do inicio ao fim e retorna o tamanho
	y := math.Abs(float64(vetor[inicio])) // Transforma em vetor[inicio] positivo
	resto := int(y) % n                   // Retorna o resto
	soma := (inicio + resto)              // Técnica para definição do índice, vai retornar um valor inteiro
	pivo_indice := soma                   // Escolha do pivô somando x da técnica a partir do ínicio

	vetor[pivo_indice], vetor[inicio] = vetor[inicio], vetor[pivo_indice] // Ocorre a troca do pivo pelo fim
	*calls++

	return particionamentoHoare(vetor, inicio, fim, calls) // Chama o método de particionar
}

func hoareRandom(vetor []int, inicio int, fim int, calls *int) {

	*calls++
	if inicio < fim {
		pivo := randomPivoHoare(vetor, inicio, fim, calls)
		hoareRandom(vetor, inicio, pivo, calls)
		hoareRandom(vetor, pivo+1, fim, calls)
	}
}

func dividir(vetorChamada []*vetorChamada) ([]*vetorChamada, []*vetorChamada) {
	return vetorChamada[0 : len(vetorChamada)/2], vetorChamada[len(vetorChamada)/2:]
}

func mergesort(vetorChamada []*vetorChamada) []*vetorChamada {
	if len(vetorChamada) <= 1 {
		return vetorChamada
	}

	left, right := dividir(vetorChamada)
	left = mergesort(left)
	right = mergesort(right)
	return intercalar(left, right)
}

func intercalar(vetorChamadaA []*vetorChamada, vetorChamadaB []*vetorChamada) []*vetorChamada {
	arr := make([]*vetorChamada, len(vetorChamadaA)+len(vetorChamadaB))

	j, k := 0, 0

	for i := 0; i < len(arr); i++ {
		if j >= len(vetorChamadaA) {
			arr[i] = vetorChamadaB[k]
			k++
			continue
		} else if k >= len(vetorChamadaB) {
			arr[i] = vetorChamadaA[j]
			j++
			continue
		}

		if vetorChamadaA[j].calls > vetorChamadaB[k].calls {
			arr[i] = vetorChamadaB[k]
			k++
		} else {
			arr[i] = vetorChamadaA[j]
			j++
		}
	}

	return arr
}
