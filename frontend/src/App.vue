<template>
  <v-app>
    <v-app-bar app color="blue-grey lighten-1" class="align-centre">
      <v-row class="justify-space-between">
        <v-col cols="1">
          <div class="d-flex">
            <a
              :href="'/favicon-' + host + '.png'"
              @click.prevent="downloadItem(item)"
            >
              <img
                :src="`favicon-${host}.ico`"
                alt="identicon"
                class="shrink mr-2"
                width="40"
              />
            </a>
          </div>
        </v-col>

        <v-col cols="3" class="d-flex justify-space-around">
          <v-menu offset-y>
            <template v-slot:activator="{ on, attrs }">
              <v-btn color="primary" dark v-bind="attrs" v-on="on">
                {{ settings.cluster }}
              </v-btn>
            </template>
            <v-list>
              <v-list-item
                v-for="(item, index) in settings.clusters"
                :key="index"
              >
                <v-list-item-title>
                  <a :href="'https://' + item">{{ item }}</a>
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>

        <v-col cols="1" style="text-align:right">
          <v-btn href="https://github.com/doddle/lander" target="_blank" text>
            <span class="mr-2"></span>
            <v-icon>mdi-open-in-new</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-app-bar>

    <v-main>
      <v-container>
        <v-row no-gutters align-content="center" justify="center">
          <template>
            <v-col v-for="obj in pieChartList" :key="obj.name" sm="3" md="3">
              <component v-bind:is="obj.content"></component>
            </v-col>
          </template>
        </v-row>
      </v-container>

      <template>
        <v-container fluid>
          <v-toolbar :color="settings.colorscheme" tabs>
            <v-tabs v-model="tab" slider-color="grey" centered>
              <v-tab v-for="item in tabList" :key="item.tab">
                {{ item.tab }}
              </v-tab>
            </v-tabs>
          </v-toolbar>
        </v-container>
      </template>
      <template>
        <v-container fluid>
          <v-tabs-items v-model="tab">
            <v-tab-item v-for="item in tabList" :key="item.tab">
              <v-card flat>
                <v-card-text>
                  <component v-bind:is="item.content"></component>
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
        </v-container>
      </template>
    </v-main>
  </v-app>
</template>

<script>
import ClusterLinks from './components/ClusterLinks'
import OverviewPieDeployments from './components/OverviewPieDeployments'
import OverviewPieNodes from './components/OverviewPieNodes'
import OverviewPieStatefulSets from './components/OverviewPieStatefulSets'
import TableNodes from './components/TableNodes'
// import axios from 'axios'

export default {
  name: 'App',
  data: function() {
    const hostLocation = location.host
    const hostName = hostLocation.split(':')[0]
    return {
      tab: null,

      tabList: [
        { tab: 'links', content: ClusterLinks },
        { tab: 'nodes', content: TableNodes }
      ],

      pieChartList: [
        { name: 'Deployments', content: OverviewPieDeployments },
        { name: 'StatefulSets', content: OverviewPieStatefulSets },
        { name: 'Nodes', content: OverviewPieNodes }
      ],
      host: hostName,
      settings: {
        colorscheme: 'blue lighten-5',
        cluster: 'unknown',
        clusters: ['cluster1.acmecorp.org']
      }
    }
  },
  title() {
    return `${this.host}`
  },

  methods: {
    async getSettings() {
      try {
        const path = '/v1/settings'
        console.log('retrieving: ' + path)
        const resp = await fetch(path)
        this.settings = await resp.json()
      } catch (error) {
        console.error(error)
      }
    }

    // downloadItem ({ url, label }) {
    //   Axios.get(url, { responseType: 'blob' })
    //       .then(response => {
    //         const blob = new Blob([response.data], { type: 'application/pdf' })
    //         const link = document.createElement('a')
    //         link.href = URL.createObjectURL(blob)
    //         link.download = label
    //         link.click()
    //         URL.revokeObjectURL(link.href)
    //       }).catch(console.error)
    // }
  },
  beforeMount() {
    this.getSettings()
  },
  components: {
    OverviewPieDeployments,
    OverviewPieStatefulSets,
    OverviewPieNodes,
    TableNodes,
    ClusterLinks
  }
}
</script>
