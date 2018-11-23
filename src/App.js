import React, { useState } from 'react'
import Map from './components/Map'
import style from './style/App/App.module.css'

function App () {
  const [markerPosition, setMarkerPosition] = useState({
    lat: 49.8419,
    lng: 24.0315
  })
  const { lat, lng } = markerPosition

  function moveMarker () {
    setMarkerPosition({
      lat: lat + 0.0001,
      lng: lng + 0.0001
    })
  }

  return (
    <div className={style.app}>
      <Map markerPosition={markerPosition} />
      <div>
        Current markerPosition: lat: {lat}, lng: {lng}
      </div>
      <button onClick={moveMarker}>Move marker</button>
    </div>
  )
}

export default App
