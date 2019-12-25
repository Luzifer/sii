const colorThresholds = {
  ok: 0,
  warn: 10,
  danger: 50,
}

const app = new Vue({

  computed: {
    dmgCargo() {
      return this.save && this.save.cargo_damage * 100 || 0
    },

    dmgTrailer() {
      return this.save && this.save.trailer_wear * 100 || 0
    },

    dmgTruck() {
      return this.save && this.save.truck_wear * 100 || 0
    },

    jobFields() {
      return [
        {
          key: 'origin_reference',
          label: 'Origin',
          sortable: true,
        },
        {
          key: 'target_reference',
          label: 'Target',
          sortable: true,
        },
        {
          key: 'cargo_reference',
          label: 'Cargo',
          sortable: true,
        },
        {
          key: 'distance',
          label: 'Dist',
          sortable: true,
        },
        {
          key: 'urgency',
          label: 'Urg',
          sortable: true,
        },
        {
          key: 'expires',
          label: 'Exp',
          sortable: true,
        },
      ]
    },

    ownedTrailers() {
      const trailers = []

      if (this.save && this.save.current_job && this.save.current_job.trailer_reference) {
        trailers.push({ value: this.save.current_job.trailer_reference, text: 'Company Trailer' })
      }

      let attachedTrailerOwned = false
      for (const id in this.save.owned_trailers) {
        if (id === this.save.attached_trailer) {
          attachedTrailerOwned = true
        }
        trailers.push({ value: id, text: this.save.owned_trailers[id] })
      }

      if (this.save && this.save.trailer_attached && !attachedTrailerOwned && (!this.save.current_job || this.save.trailer_attached !== this.save.current_job.trailer_reference)) {
        trailers.push({ value: this.save.attached_trailer, text: 'Other Trailer' })
      }

      return trailers
    },

    sortedSaves() {
      const saves = []

      for (let id in this.saves) {
        const name = this.saves[id].name !== '' ? this.saves[id].name : this.saveIDToName(id)
        saves.push({
          ...this.saves[id],
          id,
          name,
        })
      }

      return saves.sort((b, a) => { return new Date(a.save_time) - new Date(b.save_time) })
    },

    truckClass() {
      let classes = ['dmg']

      let classSelector = (prefix, value) => {
        for (let t of ['danger', 'warn', 'ok']) {
          if (value >= colorThresholds[t]) {
            return `${prefix}${t}`
          }
        }
      }

      classes.push(classSelector('truck', this.dmgTruck))

      if (this.save && this.save.trailer_attached) {
        classes.push(classSelector('trailer', this.dmgTrailer))
        classes.push(classSelector('cargo', this.dmgCargo))
      } else {
        classes.push('trailerna')
        classes.push('cargona')
      }

      return classes.join(' ')
    },
  },

  created() {
    this.loadCargo()
    this.loadProfiles()
  },

  data: {
    autoLoad: false,
    cargo: {},
    companies: {},
    jobs: [],
    profiles: {},
    save: null,
    saveLoading: false,
    saves: {},
    selectedProfile: null,
    selectedSave: null,
    showAutosaves: false,
    socket: null,
  },

  el: '#app',

  methods: {
    attachTrailer() {
      return axios.put(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/set-trailer?ref=${this.save.attached_trailer}`)
        .then(() => console.log('Trailer attached'))
        .catch((err) => console.error(err))
    },

    executeRepair(fixType) {
      return axios.put(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/fix?type=${fixType}`)
        .then(() => console.log('Repair executed'))
        .catch((err) => console.error(err))
    },

    loadCargo() {
      return axios.get(`/api/gameinfo/cargo`)
        .then((resp) => {
          this.cargo = resp.data
        })
        .catch((err) => console.error(err))
    },

    loadCompanies() {
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/companies`)
        .then((resp) => {
          this.companies = resp.data
        })
        .catch((err) => console.error(err))
    },

    loadJobs() {
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/jobs`)
        .then((resp) => {
          this.jobs = resp.data
        })
        .catch((err) => console.error(err))
    },

    loadNewestSave() {
      this.selectSave(this.sortedSaves[0].id, null)
    },

    loadProfiles() {
      return axios.get('/api/profiles')
        .then((resp) => {
          this.profiles = resp.data
        })
        .catch((err) => console.error(err))
    },

    loadSave() {
      // Load companies in background
      this.loadCompanies()

      this.saveLoading = true
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}`)
        .then((resp) => {
          this.save = resp.data
          this.saveLoading = false

          // Load jobs for that save
          this.loadJobs()
        })
        .catch((err) => console.error(err))
    },

    loadSaves() {
      if (this.socket) {
        // Dispose old socket
        this.socket.close()
        this.socket = null
      }

      let loc = window.location
      let new_uri = loc.protocol === 'https:' ? 'wss:' : 'ws:'
      new_uri += '//' + loc.host

      this.socket = new WebSocket(`${new_uri}/api/profiles/${this.selectedProfile}/saves?subscribe=true`)
      this.socket.onclose = () => this.loadSaves()
      this.socket.onmessage = evt => {
        this.saves = JSON.parse(evt.data)
        if (this.autoLoad) {
          this.loadNewestSave()
        }
      }
    },

    saveIDToName(id) {
      if (id === 'quicksave') {
        return 'Quicksave'
      }

      if (id.indexOf('autosave') >= 0) {
        return 'Autosave'
      }

      return ''
    },

    selectProfile(profileID) {
      this.selectedProfile = profileID
    },

    selectSave(saveID, evt) {
      if (this.selectedSave === saveID) {
        this.loadSave()
      } else {
        this.selectedSave = saveID
      }
    },
  },

  watch: {
    autoLoad() {
      if (this.autoLoad) {
        this.loadNewestSave()
      }
    },

    selectedProfile() {
      this.loadSaves()
    },

    selectedSave() {
      this.loadSave()
    },
  },

})
