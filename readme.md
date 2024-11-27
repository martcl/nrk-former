# NRK former - automatisk løser

NRK former _"Har du blitt hekta?"_. Nei, eller tja...

Program for å løse dagens brett, og få _"Best i Noreg"_ først! Programmet løser et vilkårlig størrelse brett `Høyde`x`Bredde` og finner løsningen med færrest mulig klikk.

## Optimaliseringer

* Brettet's state representeres med 4 unsinged integeres, en for hver farge. Hvis en farge eksisterer i en posisjon (x, y) på brettet så setter vi bit (y\*7 + x) i fargen's state til 1. Dette gjøres får å redusere minnebruk, siden A* spiser opp minne veldig kjapt. Bonus: dette gjør noen operasjoner litt kjappere, f.eks, for å sjekke om brettet er ferdig kan man sjekke med binære opperasjoner `blue_state or green_state or pink_state or orange_state == 0`.

* Distansen til mål er den naturlige logaritmen av hvor mange mulige trekk som er mulig. (ikke helt optimalt, men det beste jeg kunne finne på)

![graph_1](./assets/graph_1.png)

_(`a` er en variabel for å justere hvor "lett" man ønsker kjøre pathfinding. Hvis man setter a til noe høyt så vil den overestimere distansen og finne en løsning, men ikke den beste. Trikset er å justere a, slik at distansen blir litt underestimert hvis man ønsker å finne beste løsning)_

## Test programmet

```bash
go run cmd/main.go
```

```text
--- Board ---
3 0 0 3 0 3 2 
3 0 2 1 3 0 1 
1 0 0 2 0 1 0 
2 2 2 2 2 2 0 
1 1 0 2 1 2 3 
0 1 1 0 0 1 1 
2 1 0 2 3 0 2 
1 3 3 2 1 2 3 
0 2 2 3 2 2 3 
--------------
Iteration 0, moves: 0, esitmate: 10.033853, remainding: 37
Iteration 10000, moves: 3, esitmate: 9.523353, remainding: 31
Iteration 20000, moves: 3, esitmate: 9.615164, remainding: 32
Iteration 30000, moves: 4, esitmate: 8.654919, remainding: 23
Iteration 40000, moves: 4, esitmate: 8.779384, remainding: 24
Iteration 50000, moves: 3, esitmate: 9.790221, remainding: 34
Iteration 60000, moves: 3, esitmate: 9.873809, remainding: 35
Iteration 70000, moves: 4, esitmate: 8.898551, remainding: 25
Iteration 80000, moves: 5, esitmate: 7.932997, remainding: 18
Iteration 90000, moves: 4, esitmate: 9.012852, remainding: 26
Iteration 100000, moves: 4, esitmate: 9.012852, remainding: 26
Iteration 110000, moves: 3, esitmate: 10.033853, remainding: 37
Iteration 120000, moves: 4, esitmate: 9.122670, remainding: 27
Iteration 130000, moves: 4, esitmate: 9.122670, remainding: 27
Iteration 140000, moves: 4, esitmate: 9.122670, remainding: 27
Iteration 150000, moves: 4, esitmate: 9.122670, remainding: 27

Solution sequence, len 13:
0. (x:0, y:6)
1. (x:1, y:8)
2. (x:0, y:4)
3. (x:3, y:4)
4. (x:0, y:7)
5. (x:6, y:6)
6. (x:0, y:8)
7. (x:0, y:8)
8. (x:3, y:6)
9. (x:3, y:8)
10. (x:4, y:8)
11. (x:6, y:8)
12. (x:4, y:8)
```
