import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'

// Map outline in form of geojson points
import geoJSONmap from '../files/maps/geoJSONmap.json'

// Style
import stylecomponent from '../style/Map/Map.module.css'
import stylemap from '../style/Map/stylemap'
import styletarget from '../style/Targets/styletarget'
import stylesatellite from '../style/Satellites/stylesatellites'

// Apollo
import { ApolloClient } from 'apollo-client'
import { HttpLink } from 'apollo-link-http'
import { InMemoryCache } from 'apollo-cache-inmemory'

// Queries
import SatelliteQuery from './../queries/satelliteQuery'
import TargetsQuery from './../queries/targetsQuery'

function Map () {
  const client = new ApolloClient({
    uri: 'http://localhost:3000/graphql',
    link: new HttpLink(),
    cache: new InMemoryCache({
      addTypename: false
    })
  })

  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  useEffect(() => {
    client.query({
      query: TargetsQuery
    })
      .then(result => {
        const targetsLayer = L.geoJSON(result.data.targetFeatureCollection, {
          pointToLayer: function (feature, latlng) {
            return L.circleMarker(latlng, styletarget)
          }
        })

        mapRef.current = L.map('map', {
          zoomControl: false,
          attributionControl: false,
          zoomSnap: 0.25,
          minZoom: 2.5,
          boxZoom: true
        })
        mapRef.current.addLayer(maplayer)
        mapRef.current.addLayer(targetsLayer)
        mapRef.current.setView([0, 0], 2.5)
        mapRef.current.setMaxBounds(mapRef.current.getBounds())
      })
      .catch(error => console.error(error))
  }, [])

  useEffect(() => {
    const satelliteQuery = client.watchQuery({
      query: SatelliteQuery,
      pollInterval: 1000
    })

    const plotSats = (satellites) => {
      const satellitesLayer = L.geoJSON(satellites, {
        pointToLayer: function (feature, latlng) {
          return L.circleMarker(latlng, stylesatellite)
        }
      })

      if (mapRef.current !== null) {
        satellitesLayer.removeFrom(mapRef.current)
        mapRef.current.addLayer(satellitesLayer)
      } else {
        // mapRef.current.addLayer(satellitesLayer)
      }
    }

    satelliteQuery.subscribe({
      next: (result) => plotSats(result.data.satelliteFeatureCollection)
    })
  }, [])

  return (
    <div id='map' className={stylecomponent.map} style={stylecomponent} />
  )
}

export default Map
