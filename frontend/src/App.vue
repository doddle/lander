<template>
  <v-app>
    <v-app-bar app color="blue-grey lighten-1" class="align-centre">
      <v-row class="justify-space-between">
        <v-col cols="1">
          <div class="d-flex" v-on:click="downloadItem">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <img
                  :src="`favicon-${host}.ico`"
                  alt="identicon"
                  class="shrink mr-2"
                  width="40"
                  v-bind="attrs"
                  v-on="on"
                />
              </template>
              <span
                >This clusters unique "identicon"<br /><br />
                this is generated based on the clusters "hostname" and
                "lifecycle" (EG. prod/staging/etc)
                <br />Tip: click here to save the file to disk<br />
                ... it might be handy to have the same icons in LENs as in your
                browsers bookmarks/tabs
              </span>
            </v-tooltip>
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
                v-for="(obj, index) in settings.clusters"
                :key="index"
              >
                <v-list-item-title>
                  <a :href="'https://' + obj">{{ obj }}</a>
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

      <!-- tabs -->
      <template>
        <v-container fluid>
          <!-- a toolbar listing available and active tabs -->
          <v-toolbar :color="settings.colorscheme" tabs>
            <v-tabs v-model="activeTabName" slider-color="grey" centered>
              <v-tab
                v-for="obj in tabList"
                :key="obj.tabName"
                :disabled="isDisabled"
                @click.prevent="setActiveTabName(obj.tabName)"
              >
                {{ obj.tabName }}
              </v-tab>
            </v-tabs>
          </v-toolbar>
        </v-container>
      </template>
      <template>
        <!-- content of the tab gets templated + mounted here -->
        <v-container fluid>
          <v-tabs-items v-model="activeTabName">
            <v-tab-item v-for="item in tabList" :key="item.tab">
              <v-card flat>
                <v-card-text v-if="desiredTab === item.tabName">
                  <component v-bind:is="item.content"></component>
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
        </v-container>
      </template>
      <!-- end tabs -->
    </v-main>
  </v-app>
</template>

<script>
import ClusterLinks from './components/ClusterLinks'
import OverviewPieDeployments from './components/OverviewPieDeployments'
import OverviewPieNodes from './components/OverviewPieNodes'
import OverviewPieStatefulSets from './components/OverviewPieStatefulSets'
import TableNodes from './components/TableNodes'
import TableRoutes from './components/TableRoutes'
// import TableDeployments from './components/TableDeployments'
import axios from 'axios'

export default {
  name: 'App',
  data: function() {
    const hostLocation = location.host
    const hostName = hostLocation.split(':')[0]
    return {
      activeTabName: null,
      isDisabled: null,

      tabList: [
        { tabName: 'links', content: ClusterLinks },
        { tabName: 'nodes', content: TableNodes },
        { tabName: 'routes', content: TableRoutes }
        // { tabName: 'deployments', content: TableDeployments }
      ],
      desiredTab: 'links', // this value controls the currently rendered tab

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

  mounted() {
    this.activeTabName = this.tabList[0].tabName
  },

  methods: {
    setActiveTabName(tabName) {
      console.log('setActiveTabName: ' + tabName)
      this.desiredTab = tabName
    },
    displayContents(tabName) {
      return this.activeTabName === tabName
    },

    async getSettings() {
      try {
        const path = '/v1/settings'
        console.debug('retrieving: ' + path)
        const resp = await fetch(path)
        this.settings = await resp.json()
      } catch (error) {
        console.error(error)
      }
    },

    // clicking the cluster favicon allows a popup to save it (including name of cluster)
    downloadItem() {
      const name = 'favicon-' + this.settings.cluster + '.png'
      console.debug('downloadItem: ' + name)
      axios
        .get(name, { responseType: 'blob' })
        .then(response => {
          const blob = new Blob([response.data], { type: 'image/png' })
          const link = document.createElement('a')
          link.href = URL.createObjectURL(blob)
          link.download = name
          link.click()
          URL.revokeObjectURL(link.href)
        })
        .catch(console.error)
    }
  },
  beforeMount() {
    this.getSettings()
  },
  components: {
    OverviewPieDeployments,
    OverviewPieStatefulSets,
    OverviewPieNodes,
    TableNodes,
    // TableDeployments,
    ClusterLinks
  }
}
</script>
