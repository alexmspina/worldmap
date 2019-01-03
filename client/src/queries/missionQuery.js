import gql from 'graphql-tag'

const MissionQuery = gql`{
    satelliteFeatureCollection{
        features {
            properties {
                id,
                velocity,
                altitude,
                mission {
                    id,
                    config,
                    gatewayID,
                    gatewayOBAnt,
                    gatewayMaxPointingTime,
                    beams {
                        id,
                        epcs,
                        targetOBAnt,
                        targetMaxPointingTime,
                        camp,
                        campMode,
                        campGain,
                        ldla,
                        ldlaMode,
                        ldlaFcaGain,
                        ldlaGcaGain,
                        ldlaScaGain
                    }
                },
            },
        },
    }
}`

export default MissionQuery
