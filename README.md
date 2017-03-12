# TTK4145-Spring-2017-Group-89

Good GO-tutorial:

https://www.youtube.com/watch?v=CF9S4QZuV30&t=923s


# TODO:
* Finne ut hva watchdog skal gjøre
* Controller-flyt for knappetrykk, ender med å sende over nettverk.
* Recovery fra storage (Hva skal master og controller gjøre når det finnes stuff i storage?)
* Fikse channels for å unngå farlig delte vaiabler i "slave"
  * Innebærer både ordre-lista og "elevatorState"


# Bugs:
* Dersom master.updateOrders() får inn samme etasje som heisen er på for øyeblikket blir orderen plassert på feil sted.

* Nettverk er vanskelig å få til å fungere. Saken er at om jeg ikke closer en socket/connect(whatever it's called) så får jeg problemer neste gang jeg kjører programmet. Tror kanskje det kan være lurt å definere og close sockets helt i main.go, men vi må skrive om flere funksjoner da.
  * Er til en viss grad fikset (12/03/17)
