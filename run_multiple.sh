for ants in 20 50 100 200 500 
do
    for gens in 20 50 100 200 500
    do
        for evap in 0.1 0.2 0.3 0.5 
        do
            for data in datasets/easy.col datasets/hard.col
            do 
                go run . -ants $ants -gens $gens -minpheromone 1 -maxpheromone 10 -evaporation $evap -file $data -seed 4132 -getstats 
            done
        done
    done
done