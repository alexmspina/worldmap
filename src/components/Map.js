import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'
import stylecomponent from '../style/Map/Map.module.css'
import geoJSONmap from '../files/maps/geoJSONmap.json'
import stylemap from '../style/Map/stylemap'
// import targetsFile from '../files/targets/TARGETS_275.json'
import GetTargets from './TargetsQuery'
import { Query } from 'react-apollo'
import styletarget from '../style/Targets/styletarget'
// import { GeoJSON } from 'react-leaflet'
// import { ApolloClient } from 'apollo-client'

function Map () {
  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  useEffect(() => {
    mapRef.current = L.map('map', {
      zoomControl: false,
      attributionControl: false,
      zoomSnap: 0.25,
      minZoom: 2.5,
      boxZoom: true
    })

    mapRef.current.addLayer(maplayer)
    mapRef.current.setView([0, 0], 2.5)
    mapRef.current.setMaxBounds(mapRef.current.getBounds())
  }, [])

  return (
    <div id='map' className={stylecomponent.map} style={stylecomponent}>
      <Query query={GetTargets}>
        {({ loading, error, data }) => {
          if (loading) return 'Loading...'
          if (error) return `Error! ${error.message}`

          console.log(data.targetFeatureCollection)

          const targetsLayer = L.geoJSON(data.targetFeatureCollection, {
            pointToLayer: function (feature, latlng) {
              return L.circleMarker(latlng, styletarget)
            }
          })

          mapRef.current.addLayer(targetsLayer)

          return <div />
        }}
      </Query>
    </div>
  )
}

export default Map
