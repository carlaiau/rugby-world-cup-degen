# 🏉🏉🏉 Rugby World Cup 2023 Degen 🏉🏉🏉

A simple Golang and React hack project that saved me looking at multiple bookies and odds aggregators when line shopping.

If you simply want a punt, you can use this to find the best lines.
Alternatively you could use a [middle betting strategy](https://www.bestodds.com/guides/middle-bets/) which is why there is color guidance and emphasis on the Δ between handicaps.

The Server pulls odds from the-odds-api.com and the NZ TAB private API every 5 minutes if the client has pinged within the last minute. 
Client will keep pinging the server every minute whilst the react app is running.

### Screenshot
<img width="600" alt="Screenshot 2023-10-17 at 12 28 33 AM" src="https://github.com/carlaiau/rugby-world-cup-degen/assets/6896663/0b806db8-37f0-4c63-bcfb-a56f0c4fd3ca">
<img width="600" alt="Screenshot 2023-10-17 at 12 28 48 AM" src="https://github.com/carlaiau/rugby-world-cup-degen/assets/6896663/6ffe550a-11f4-495f-aa30-3b4cdfc55849">

### Expectations:

- you create a environment file at `server/.env` with the contents: `THE_ODDS_API_KEY="someKeyThatYouCopyPasteFromThere"`.
- You're in NZ or Australia or the NZ TAB private API will block your requests.
- All Blacks win the tournament.

## Recipe

### Run Server
1) cd server
2) go build
3) ./degen

### Run Client
1) cd client
2) npm install
3) npm start
4) visit localhost:8000

### Notes
I didn't want to burn time deploying this remotely but it is trivial to access on mobile if you change the network address at the top of `client/src/App.js` to your local network address.

At the top of `client/src/App.js` there is a dict which indicates which bookmaker you have accounts with along with the betting page, of said bookmaker so you can amend this to add further bookies. 

This object dictates the green rows in the app (shown in screenshot) indicating which bookies you can actually use. for instance it could be modified to
```
const SPORTSBOOK_WITH_ACCOUNT = {
  "TopSport": "https://www.topsport.com.au/Sport/Rugby_Union/Rugby_World_Cup_Matches/Matches",
  "Pinnacle": "https://www.pinnacle.com/en/rugby-union/rugby-world-cup/matchups/",
}
```
And then only TopSport and Pinnacle would show up green. 

The Delta is based on all available markets, irrespective of where you have an account.

No support will be given to this project, happy punting.





