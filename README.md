# TTK4145-Spring-2017-Group-89

Good GO-tutorial:

https://www.youtube.com/watch?v=CF9S4QZuV30&t=923s


# TODO:
* Master skal lytte jevnlig etter JSON-objekter sendt fra slaves.
* Slave skal tilsvarende sende JSON-objekt
* Master skal sende ordre lagret i JSON-objekt
* Finne ut hva watchdog skal gjøre
* Skrive algoritme for å bestemme hvilken heis som skal ta jobben (master-funksjonalitet)
* Dekode mottat melding fra JSON til struct -> Lagre
* Master -> Controller for master-PCen (muligens localhost)
* Controller-flyt for knappetrykk, ender med å sende over nettverk.
* Egen funksjon for å oppdatere elevatorState (modifikasjon av PrintLastFloorIfChanged)
* Fjerne kommandoer fra order-liste
* Lage id for sletting av ordre (inkrementell eller random?)
* Recovery fra storage (Hva skal master og controller gjøre når det finnes stuff i storage?)
* Nummerering av slaves. Heis nr. ????

# Bugs:
* Dersom master.updateOrders() får inn samme etasje som heisen er på for øyeblikket blir orderen plassert på feil sted.

* Nettverk er vanskelig å få til å fungere. Saken er at om jeg ikke closer en socket/connect(whatever it's called) så får jeg problemer neste gang jeg kjører programmet. Tror kanskje det kan være lurt å definere og close sockets helt i main.go, men vi må skrive om flere funksjoner da.