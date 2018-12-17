import React, { useEffect, useRef } from 'react'
import * as L from 'leaflet'
import stylecomponent from '../style/Map/Map.module.css'
import geoJSONmap from '../files/maps/geoJSONmap.json'
import stylemap from '../style/Map/stylemap'
import styletarget from '../style/Targets/styletarget'
import targetsFile from '../files/targets/TARGETS_275.json'
import gql from 'graphql-tag'
// import { graphql } from 'react-apollo'
import { SubscriptionClient } from 'subscriptions-transport-ws'
import { ApolloClient } from 'apollo-boost'
import { WebSocketLink, HttpLink } from 'apollo-link-ws'
import { split } from 'apollo-link'
import { getMainDefinition } from 'apollo-utilities'

const httpLink = new HttpLink({
  uri: 'http://localhost:3000/targets'
})

const wsLink = new WebSocketLink({
  uri: `ws://localhost:3000/subscriptions`,
  options: {
    reconnect: true
  }
})

const link = split(
  ({ query }) => 
)

// const GRAPHQL_ENDPOINT = 'ws://localhost:8080/subscriptions'

// const client = new SubscriptionClient(GRAPHQL_ENDPOINT, {
//   reconnect: true
// })

// const apolloClient = new ApolloClient({
//   networkInterface: client
// })

// console.log('just before the query')

// apolloClient
//   .query({
//     query: gql`
//       {
//         satellites
//       }
//     `
//   })
//   .then(result => console.log(result))

function Map ({ svgMapBounds }) {
  // create map reference
  const mapRef = useRef(null)

  // create map from geoJSON layer
  const maplayer = L.geoJSON(geoJSONmap, {
    style: stylemap
  })

  // create target layer from geoJSON files
  const targetsLayer = L.geoJSON(targetsFile, {
    pointToLayer: function (feature, latlng) {
      return L.circleMarker(latlng, styletarget)
    }
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
    mapRef.current.addLayer(targetsLayer)
    mapRef.current.setView([0, 0], 2.5)
    mapRef.current.setMaxBounds(mapRef.current.getBounds())
  }, [])

  return <div id='map' className={stylecomponent.map} style={stylecomponent} />
}

export default Map
