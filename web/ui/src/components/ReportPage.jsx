import React, { Component } from 'react'
import { Pane } from 'evergreen-ui'
import { Provider, Subscribe } from 'unstated'

import ReportList from './ReportList'

import ReportContainer from '../containers/ReportContainer'

const reportContainer = new ReportContainer()

export default class ReportsPage extends Component {
  render () {
    return (
      <Provider inject={[reportContainer]}>
        <Subscribe to={[ReportContainer]}>
          {(reportStore) => (
            <Pane>
              <ReportList reportStore={reportStore} />
            </Pane >
          )}
        </Subscribe>
      </Provider >
    )
  }
}
