package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// Struct do conteiner
type conteiner struct {
	ordem           int
	codigo          string
	cnpj            string
	peso            string
	cnpjErrado      string
	porcentagemPeso int
	desvioPeso      int
}

func main() {

	//start := time.Now()

	var conteineres []*conteiner
	var conteineresCnpjErrado []*conteiner
	var conteineresPesoErrado []*conteiner

	arquivoEntrada := os.Args[1] // Entrada
	arquivoSaida := os.Args[2]   // Saída

	// Leitura de arquivo
	arquivo, err := os.Open(arquivoEntrada)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo")
		return
	}
	scanner := bufio.NewScanner(arquivo)

	// Listagem de todos os conteineres no porto
	if scanner.Scan() {
		qtdConteineres, _ := strconv.Atoi(scanner.Text())

		for i := 0; i < qtdConteineres && scanner.Scan(); i++ { // n
			linha := scanner.Text()
			palavra := strings.Fields(linha)
			conteineres = append(conteineres, &conteiner{i, palavra[0], palavra[1], palavra[2], "", 0, 0})
		}
	}

	conteineresOrdenadosCodigo := mergesortCodigo(conteineres) // n log n

	// Listagem de todos conteineres fiscalizados
	if scanner.Scan() {
		qtdConteineresFiscalizados, _ := strconv.Atoi(scanner.Text())
		for i := 0; i < qtdConteineresFiscalizados && scanner.Scan(); i++ { // n log n
			linha := scanner.Text()
			palavra := strings.Fields(linha)

			problemaCnpj := false

			fiscalizadoInterface := buscaBinaria(conteineresOrdenadosCodigo, palavra[0])
			fiscalizado, _ := fiscalizadoInterface.(*conteiner)

			if fiscalizado.cnpj != palavra[1] {
				conteineresCnpjErrado = append(conteineresCnpjErrado, &conteiner{fiscalizado.ordem, fiscalizado.codigo, fiscalizado.cnpj, fiscalizado.peso, palavra[1], 0, 0})
				problemaCnpj = true
			}

			pesoConteiner, _ := strconv.ParseFloat(fiscalizado.peso, 32)
			pesoConteinerFiscalizado, _ := strconv.ParseFloat(palavra[2], 32)

			desvioPesoAux := math.Abs(pesoConteiner - pesoConteinerFiscalizado)

			consultaPeso := Round(math.Abs((pesoConteiner-pesoConteinerFiscalizado)/pesoConteiner) * 100)

			resultadoConsulta := int(consultaPeso)

			desvioPeso := int(desvioPesoAux)

			if resultadoConsulta > 10 && !problemaCnpj {
				conteineresPesoErrado = append(conteineresPesoErrado, &conteiner{fiscalizado.ordem, fiscalizado.codigo, fiscalizado.cnpj, fiscalizado.peso, "", resultadoConsulta, desvioPeso})
			}
		}

		conteineresCnpjErradoOrdenado := mergesort(conteineresCnpjErrado) // n log n

		conteineresPesoErradoOrdenado := mergesort(conteineresPesoErrado)                       // n log n
		conteineresPesoErradoOrdenado = mergesortPesoPorcentagem(conteineresPesoErradoOrdenado) // n log n

		f, errOut := os.Create(arquivoSaida)
		if errOut == nil {
			w := bufio.NewWriter(f)
			for _, v := range conteineresCnpjErradoOrdenado {
				w.WriteString(fmt.Sprintf("%s: %s<->%s\r\n", v.codigo, v.cnpj, v.cnpjErrado))
			}

			for _, v := range conteineresPesoErradoOrdenado {
				w.WriteString(fmt.Sprintf("%s: %dkg (%0.0d%%)\r\n", v.codigo, v.desvioPeso, v.porcentagemPeso))
			}

			w.Flush()
			f.Close()
		}

	}

	arquivo.Close()
	//timeTrack(start, "TEMPO")

}

func Round(f float64) float64 {
	if f < 0 {
		return math.Ceil(f - 0.5)
	}
	return math.Floor(f + 0.5)
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s => %s", name, elapsed)
}

func buscaBinaria(conteiner []*conteiner, valor string) interface{} {

	inicio := 0
	fim := len(conteiner) - 1
	for inicio <= fim {
		media := (inicio + fim) / 2
		if strings.Compare(valor, conteiner[media].codigo) == 0 {
			return conteiner[media]
		} else if strings.Compare(valor, conteiner[media].codigo) > 0 {
			inicio = media + 1
		} else {
			fim = media - 1
		}
	}

	if inicio == len(conteiner) {
		return -1
	} else {
		return 0
	}
}

func dividir(conteiner []*conteiner) ([]*conteiner, []*conteiner) {
	return conteiner[0 : len(conteiner)/2], conteiner[len(conteiner)/2:]
}

func mergesortCodigo(conteiner []*conteiner) []*conteiner {
	if len(conteiner) <= 1 {
		return conteiner
	}
	left, right := dividir(conteiner)
	left = mergesortCodigo(left)
	right = mergesortCodigo(right)
	return intercalarCodigo(left, right)
}

func mergesortPesoPorcentagem(conteiner []*conteiner) []*conteiner {
	if len(conteiner) <= 1 {
		return conteiner
	}
	left, right := dividir(conteiner)
	left = mergesortPesoPorcentagem(left)
	right = mergesortPesoPorcentagem(right)
	return intercalarPorcentagemPeso(left, right)
}

func mergesort(conteiner []*conteiner) []*conteiner {
	if len(conteiner) <= 1 {
		return conteiner
	}

	left, right := dividir(conteiner)
	left = mergesort(left)
	right = mergesort(right)
	return intercalar(left, right)
}

func intercalarCodigo(conteinerA, conteinerB []*conteiner) []*conteiner {
	arr := make([]*conteiner, len(conteinerA)+len(conteinerB))

	j, k := 0, 0

	for i := 0; i < len(arr); i++ {
		if j >= len(conteinerA) {
			arr[i] = conteinerB[k]
			k++
			continue
		} else if k >= len(conteinerB) {
			arr[i] = conteinerA[j]
			j++
			continue
		}

		if conteinerA[j].codigo > conteinerB[k].codigo {
			arr[i] = conteinerB[k]
			k++
		} else {
			arr[i] = conteinerA[j]
			j++
		}
	}

	return arr
}

func intercalarPorcentagemPeso(conteinerA, conteinerB []*conteiner) []*conteiner {
	arr := make([]*conteiner, len(conteinerA)+len(conteinerB))

	j, k := 0, 0

	for i := 0; i < len(arr); i++ {
		if j >= len(conteinerA) {
			arr[i] = conteinerB[k]
			k++
			continue
		} else if k >= len(conteinerB) {
			arr[i] = conteinerA[j]
			j++
			continue
		}

		if conteinerA[j].porcentagemPeso < conteinerB[k].porcentagemPeso {
			arr[i] = conteinerB[k]
			k++
		} else {
			arr[i] = conteinerA[j]
			j++
		}
	}

	return arr
}

func intercalar(conteinerA, conteinerB []*conteiner) []*conteiner {
	arr := make([]*conteiner, len(conteinerA)+len(conteinerB))

	j, k := 0, 0

	for i := 0; i < len(arr); i++ {
		if j >= len(conteinerA) {
			arr[i] = conteinerB[k]
			k++
			continue
		} else if k >= len(conteinerB) {
			arr[i] = conteinerA[j]
			j++
			continue
		}

		if conteinerA[j].ordem > conteinerB[k].ordem {
			arr[i] = conteinerB[k]
			k++
		} else {
			arr[i] = conteinerA[j]
			j++
		}
	}

	return arr
}
