package NodeSingleConnection

/*[FFF]
1. Inneholder en tråd som skal bli kalt hver gang en ny Node vil koble seg på.
2. Denne vil motta data den skal sende over en channel hvor den andre enden ligger lagret i NodeRoutingTable listen sammen med all den andre data om hvilken node som denne tråden kobler til.
3. Har i oppdrag å pinge og sørge for at forbindelsen opprettholdes.
4. Har i oppdrag å motta meldinger fra annen node og sende til NodeMessageRelay
5. Har i oppdrag å sende meldinger til annen node som er mottatt fra NodeMessageRelay 
6. Må varsle NodeConnectionManager hvis den mister forbindelsen/timer ut.
(7.) Føre statestikk over hvor mye som sendes.
*/

import
(
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	
	/* Implement the notes above */
}