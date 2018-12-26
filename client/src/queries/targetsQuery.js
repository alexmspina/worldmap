import gql from 'graphql-tag'

const TargetsQuery = gql`{
  targetFeatureCollection{
    type,
    features {
      type,
      geometry {
        type,
        coordinates,
      },
      properties {
        gatewayFlag,
        shortName,
      },
    }
  }
}`

export default TargetsQuery
