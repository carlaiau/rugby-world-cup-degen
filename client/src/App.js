import axios from "axios";
import { useEffect, useState } from "react";
import { format } from 'date-fns';
import { formatDistance } from 'date-fns';
import './App.css';
const API_ADDRESS = "http://localhost:8080/"

const SPORTSBOOK_WITH_ACCOUNT = {
  "NZ TAB": "https://www.tab.co.nz/sports/competition/18525/rugby-union/international/world-cup-2023/matches",
  "SportsBet": "https://www.sportsbet.com.au/betting/rugby-union/world-cup-2023",
  "TopSport": "https://www.topsport.com.au/Sport/Rugby_Union/Rugby_World_Cup_Matches/Matches",
  "Pinnacle": "https://www.pinnacle.com/en/rugby-union/rugby-world-cup/matchups/",
  "Unibet": "https://www.unibet.com.au/betting/sports/filter/rugby_union/rugby_world_cup_2023/all/matches"
}

const renderMarkets = (event, key) => {
  const sortedMarkets = event.marketsByName[key].
  sort((a, b) => a.bookmaker.localeCompare(b.bookmaker)).
  sort((a, b) => a.line - b.line)
  if(!sortedMarkets.length)
  return
  const bestSpreadDifference = sortedMarkets[sortedMarkets.length - 1].line - sortedMarkets[0].line

  const spreadClass = bestSpreadDifference > 3.5 ? 'good' : bestSpreadDifference > 1.5 ? 'okay': ''

  return (
    <>
      <thead>
        <tr>
          <th colSpan={2}>
            <div className="title">
              {key}
              <span className="time">updated at {format(new Date(sortedMarkets[0].last_update), 'hh:mm')} / {format(new Date(sortedMarkets[sortedMarkets.length-1].last_update), 'hh:mm')}</span>
            </div>
          </th>
          <th colSpan={1} className={'spread ' + spreadClass}>Δ {bestSpreadDifference}</th>
        </tr>
        { key != "Total" ?
          <tr>
            <th>Bookie</th>
            <th>Handicap</th>
            
            <th>Odds</th>
            {
              /*
              <th>Juice</th>
              <th>Updated</th>
              */
            }

          </tr>
          : 
          <></>
        }
        
      </thead>
      <tbody>
      {sortedMarkets.map(
        (market) => {

          const link = SPORTSBOOK_WITH_ACCOUNT[market.bookmaker]
          let juice = 0
          market.outcomes.forEach((outcome) => {
            juice += 1/outcome.price
          })
          return (
            <tr className={link ? 'with-account': ""}>
              <td>
                { link ? 
                <a href={link} target="_blank">
                  {market.bookmaker}
                </a>
                  :
                  market.bookmaker
                }
              </td>
              
              <td>{market.line.toFixed(1)}</td>
              
              <td className="odds">
              {
                market.outcomes.map((outcome, index) => {
                  return (
                    <span>
                      {outcome.price.toFixed(2)}
                      {index == 0 ? '/' : ''}
                     
                    </span>
                    
                  )
                })
              }
              </td>
              
            </tr>
          )
        })}
      </tbody>
    </>
  )

}
function App() {

  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);
  useEffect(() => {
    

    fetchData();
    const intervalId = setInterval(() => {
      
      fetchData()
    }, 60000);
    return () => clearInterval(intervalId);
    

  }, []);


  const fetchData = async () => {
    const date = new Date()
    console.log("Last Fetch", date)
    setLoading(true)
    const response = await axios.get(API_ADDRESS);
    setLoading(false)
    setData(response.data)
  };
  

  if(data.length == 0){

    return (
      <div className="App">
        <p>Empty</p>
      </div>
    );
  }


  // do the manipulation here

  data.forEach((event) => {
    event.marketsByName = {
      'Line': event.markets.filter((market) => market.name === 'LINE' && market.outcomes.length === 2),
      'Total': event.markets.filter((market) => market.name === 'TOTAL' && market.outcomes.length === 2),
    }
  })

  // determine if we need to flip flop the TAB's handicap
  data.forEach((event) => {

    
    let negativeCount = 0
    let positiveCount = 0 
    
    event.marketsByName['Line'].forEach((market) => {
      if (market.bookmaker === 'NZ TAB') {
        return
      }
      if (market.line < 0) {
        negativeCount++
      }
      else if(market.line > 0) {
        positiveCount++
      
      }
    })
    if (negativeCount > positiveCount) {
      event.marketsByName['Line'].forEach((market) => {
        if (market.bookmaker === 'NZ TAB' && market.line > 0) {
          market.line = -market.line
        }
      })
    }
    else if (negativeCount < positiveCount) {
      event.marketsByName['Line'].forEach((market) => {
        if (market.bookmaker === 'NZ TAB' && market.line < 0) {
          market.line = -market.line
        }
      })
    }

    else {
      // do nothing
    }

    // 
  })
  
  return (

    <div className="app">
    
    {
    data.sort((a, b) => {
      const timeA = new Date(a.startTime).getTime();
      const timeB = new Date(b.startTime).getTime();
      return timeA - timeB;
    }).map((event, index) => {

      const startTime =  new Date(event.startTime).getTime();
      const currentTime = new Date().getTime();
      let isLive = false
      const timeUntil = currentTime - startTime;
      if(startTime < currentTime){
        isLive = true
      }
      return(
        <div className={isLive ? 'container is-live' : 'container'}>
          <div className="title-container">
            <h3 key={event.name}>{event.name}</h3>
            <h3 className="date">{isLive ? 
              <span className="is-live">Live</span>
              : <></>
              }{
                format(new Date(event.startTime), 'EEE, dd MMM HH:mm')}
                <br/>
                {formatDistance(startTime, currentTime, {includeSuffix: true})}
              
            </h3>
          </div>
          
          <table style={{ width: "100%" }}>
          {
            renderMarkets(event, "Line")
          }
          {
            renderMarkets(event, "Total")
          }
          </table>
        </div>
      )
    })
    }
    </div>
  )
}

export default App;
