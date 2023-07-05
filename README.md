# AOC - Clique Máximo

## Francisco Bonome Andrade & Lucas Starling de Paula Salles


O trabalho foi desenvolvido [nesse repositório do GitHub](https://github.com/franciscobonand/aoc-max-clique).

O Problema do Clique Máximo é um problema clássico da teoria dos grafos e da complexidade computacional.
Dado um grafo não direcionado, o objetivo é encontrar o maior clique, que é um subconjunto de vértices em que cada par de vértices está conectado por uma aresta.
Encontrar o clique máximo em um grafo é um problema NP-difícil, o que significa que nenhum algoritmo conhecido pode resolvê-lo em tempo polinomial para todas as instâncias.
Portanto, abordagens alternativas, como métodos heurísticos, tornam-se essenciais para aplicações práticas.  

A Otimização de Colônia de Formigas (ACO) é uma meta-heurística inspirada no comportamento de busca por alimento das formigas.
É uma técnica de otimização poderosa que tem sido aplicada com sucesso em vários problemas combinatórios.
A ideia central do ACO é simular o comportamento das formigas em busca de alimento,
construindo e melhorando iterativamente soluções candidatas por meio de trilhas de feromônio e heurísticas de busca local.  

O objetivo deste projeto é implementar um algoritmo de ACO especialmente adaptado para enfrentar o Problema do Clique Máximo.
Aproveitando as vantagens inerentes da otimização de colônia de formigas,
pretendemos desenvolver uma solução eficiente e eficaz que possa encontrar cliques próximos ao ótimo ou ótimos dentro de um grafo dado.  

Ao longo desta documentação, discutiremos os fundamentos teóricos do Problema do Clique Máximo,
forneceremos uma visão geral do algoritmo de Otimização de Colônia de Formigas,
descreveremos os detalhes de projeto e implementação de nossa solução e apresentaremos os resultados de nossos experimentos.  

## Como executar o programa

Primeiramente faça o download e instale a [versão mais recente da linguagem Golang](https://go.dev/doc/install).  
Com a instalação realizada, basta executar o seguinte comando do diretório raiz desse repositório:

```sh
go run .
```

**Caso não deseje instalar o Golang, pode optar por executar o binário que se encontra na pasta `/bin`**  
Para isso, primeiro execute o comando:

```sh 
chmod +x ./bin/aoc-max-clique
```

E então execute o programa com:

```sh 
./bin/aoc-max-clique
```

### Flags - Parametrização

Ao executar o programa, flags podem ser utilizadas para definir alguns parâmetros da execução.  
As flags são definidas da forma `.bin/aoc-max-clique -flag1 valor1 -flag2 valor2` (ou `go run . -flag1 valor1 -flag2 valor2`)  
São elas:

| Flag         | Default | Tipo          | Descrição                               |
| ------------ | ------- | ------------- | --------------------------------------- |
| ants         | 100     | Int           | Número de formigas                      |
| gens         | 10      | Int           | Quantidade de gerações/iterações        |
| minpheromone | 0.1     | Float         | Valor mínimo de feromônio em um vértice |
| maxpheromone | 5       | Float         | Valor máximo de feromônio em um vértice |
| evaporation  | 0.2     | 0 < Float < 1 | Taxa de evaporação de feromônios        |
| file         | "easy"  | String        | Nome do arquivo de entrada              |
| seed         | 1       | Int           | Semente aleatória                       |
| verbose      | false   | Bool          | Imprime estatísticas por geração        |
| getstats     | false   | Bool          | Gera relatório de execução              |

Exemplo:

```sh
go run . -ants 20 -gens 20 -minpheromone 0.1 -maxpheromone 10 -evaporation 0.1 -file "datasets/easy.col" -seed 4132
```

Também é possível ver a descrição das flags usando `--help`:

```sh
go run . --help
// OU
./bin/aoc-max-clique --help
```

## Implementação

Nesse tópico serão apresentadas as principais estruturas utilizadas no programa, assim como decisões de implementação e limitações.  

### Representação do problema

Uma das principais decisões a serem feitas foi a de como representar o grafo, assim como quais componentes desse grafo irão receber feromônios.  
Em nossa implementação, o grafo a ser analisado é representado por uma matriz de adjacência, e cada célula (i, j) da matriz contém um valor numérico.
Esse valor representa a quantidade de feromônio que existe na aresta que liga o ponto i ao j:

![Representação do problema](/images/matrix.svg "Matriz de adjacência, com feromônios nas arestas")

### Construção da solução

No ACO, para cada iteração/geração, deve-se construir uma solução para cada uma das formigas.
Para o problema do Clique Máximo, portanto, esse passo consiste em construir um clique no grafo, e isso é feito pela função `buildClique`.  
Para cada formiga, o passo-a-passo dessa função dá-se da seguinte maneira:  

1.  É escolhido um vértice aleatório do grafo como inicial, e ele é adicionado ao conjunto clique.
2.  Cria-se um conjunto de membros candidatos a entrarem no clique, que consiste nos vértices vizinhos ao vértice inicial.
3.  Enquanto o conjunto de candidatos não estiver vazio, realiza-se os seguintes passos para adicionar um novo vértice ao clique:
    1.  Obtem-se o fator de feromônio para cada candidato que "deseja" entrar no clique.
    O fator de cada candidato consiste na soma dos feromônios das arestas que ligam esse vértice ao clique já existente, elevado a 2.
    Dessa forma todos os vértices candidatos podem, inicialmente, serem adicionados ao clique, porém aqueles que possuem mais arestas (ou arestas com mais feromônios) que o conectam ao clique tem fator maior.
    2.  Escolhe-se o candidato a ser adicionado ao clique de maneira aleatória, levando em conta o "peso" de cada candidato com base no valor dos fatores previamente calculados.
    3.  O candidato selecionado é adicionado ao clique, e a lista de candidatos é atualizada.
    **Essa atualização consiste na interseção entre os candidatos anteriores com os vizinhos do vértice que foi adicionado ao clique**.
    Dessa forma, garante-se que a solução encontrada sempre será um clique.

### Atualização dos feromônios

Após a construção das soluções para cada formiga, deve-se atualizar as quantidades dos feromônios de todas as arestas do grafo.  
Nessa implementação, optou-se por utilizar o método elitista para atualização dos feromônios.
Ou seja, primeiramente todos os feromônios do grafo são evaporados conforme o valor da taxa de evaporação e, em seguida, os vértices que compreendem a melhor solução encontrada pelo caminho de uma formiga recebem mais feromônios.  
O passo-a-passo da função que implementa a evaporação dá-se da seguinte maneira:  

1.  Dadas as soluções encontradas no passo anterior, escolhe-se a melhor (maior clique).
    1.  Se essa solução for a melhor já encontrada entre todas as iterações prévias, atualiza o valor do melhor clique global já encontrado.
2.  Para todas as arestas do grafo, decai o valor do feromônio com base na taxa de evaporação.
Esse novo valor nunca é menor que o valor mínimo de feromônio que uma aresta pode possuir (parâmetro definido no início da execução do programa).
3.  Para todas as arestas que pertencem ao melhor clique encontrado nessa iteração, aumenta-se o valor de feromônios.
Esse novo valor nunca é maior que o valor máximo de feromônio que uma aresta pode possuir (parâmetro definido no início da execução do programa).
O valor adicionado é definido por `Quantidade Atual + (1 / (1 + (Melhor Solução Global - Melhor Solução da Iteração)))`.
    1.  Caso a Melhor Solução da Iteração seja a Melhor Solução Global, a quantidade atual de feromônios é acrescida em 1.
    2. Caso a Melhor Solução Global seja maior, isso implica que a quantidade de feromônios adicionada é inferior a 1.
    Isso é feito com o intuito de permitir a encontrabilidade de novas soluções, porém favorece o *exploitation* de soluções melhores previamente encontradas.

## Análises

As análises realizadas a partir dos [dados fornecidos](/datasets) podem ser encontradas no [Jupyter Notebook presente nesse repositório](CompNatTP2.ipynb).  
Os resultados usados para análise podem ser encontrados na pasta ['analysis'](/analysis).

### Os Problemas

Foram disponibilizados três conjuntos de dados representando três problemas grafos de grande dimensionalidade:

- Easy: Grafo com 500 nós, 62624 arestas, contendo clique maximo de tamanho 13.
- Hard: Grafo com 700 nós, 121728 arestas, contendo clique maximo de tamanho 44.
- Harder: Grafo com 600 nós, 207643 arestas, contendo clique maximo de tamanho 26. Esse Problema se mostrou mais simples que o 'Hard' apesar de conter um grafo de maior dimensão.

### Parametros de otimização

Para o processo de otimização por colonia de formigas existem três parametros relevantes à serem analisados:
- Número de formigas - *ants*: Contagem do número de formigas utilizadas

Analisaremos a performance do algoritmo de otimização por colonia de formigas implementado separadamente para os três problemas, isso porque observamos comportamentos divergentes com relação




## Conclusão

A implementação da otimização por colônia de formigas (ACO) para resolver o problema do máximo clique apresentou resultados promissores em diversos cenários. Ao lidar com grafos que possuem uma solução simples, a implementação do ACO foi capaz de encontrar o máximo clique facilmente e de forma eficiente. A capacidade do algoritmo de convergir rapidamente para uma solução é especialmente notável, uma vez que ele identificou diferentes máximos cliques em cada execução quando múltiplas soluções existiam.  

No entanto, ao enfrentar instâncias mais desafiadoras do problema do máximo clique, a implementação do ACO teve dificuldade em encontrar a solução ótima de forma consistente. Apesar de ajustar parâmetros como a quantidade de formigas, o número de iterações e a taxa de evaporação do feromônio, o desempenho do algoritmo ficou aquém do esperado. Essa limitação provavelmente se deve à abordagem elitista utilizada na atualização dos feromônios.  

Para melhorar a eficácia da implementação do ACO na resolução de instâncias mais difíceis do problema do máximo clique, estratégias alternativas podem ser exploradas. Por exemplo, investigar diferentes formas de atualizar os feromônios que priorizem a *exploration* em detrimento da *exploitation* pode levar a melhores resultados. Além disso, incorporar heurísticas de busca local ou técnicas de diversificação pode aprimorar a capacidade do algoritmo de escapar de ótimos locais e explorar mais abrangentemente o espaço de soluções.  

Em suma, a implementação da otimização por colônia de formigas para o problema do máximo clique apresenta uma abordagem promissora com grande potencial. Embora o algoritmo tenha se destacado em casos mais simples e tenha demonstrado uma convergência eficiente, é necessário realizar refinamentos adicionais para superar os desafios apresentados por instâncias mais complexas. Ao abordar as limitações e explorar técnicas alternativas, futuras iterações dessa implementação do ACO podem buscar alcançar soluções mais ótimas para uma variedade maior de problemas de máximo clique.  
