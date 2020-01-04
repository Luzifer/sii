const colorThresholds = {
  ok: 0,
  warn: 10,
  danger: 50,
}

window.app = new Vue({

  computed: {
    cargoSelectItems() {
      const result = []

      for (const ref in this.cargo) {
        result.push({
          value: ref,
          text: `${this.cargo[ref].name}`,
        })
      }

      return result.sort((a, b) => a.text.localeCompare(b.text))
    },

    companySelectItems() {
      const result = []

      for (const ref in this.companies) {
        result.push({
          value: ref,
          text: this.companyNameFromReference(ref),
        })
      }

      return result.sort((a, b) => a.text.localeCompare(b.text))
    },

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

      for (const id in this.saves) {
        const name = this.saves[id].name !== '' ? this.saves[id].name : this.saveIDToName(id)
        saves.push({
          ...this.saves[id],
          id,
          name,
        })
      }

      return saves.sort((b, a) => new Date(a.save_time) - new Date(b.save_time))
    },

    truckClass() {
      const classes = ['dmg']

      const classSelector = (prefix, value) => {
        for (const t of ['danger', 'warn', 'ok']) {
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
    newJob: { weight: 10 },
    plannedRoute: [],
    profiles: {},
    save: null,
    saveLoading: false,
    saves: {},
    selectedProfile: null,
    selectedSave: null,
    showAutosaves: false,
    showSaveModal: false,
    socket: null,
  },

  el: '#app',

  methods: {
    attachTrailer() {
      this.showSaveModal = true
      return axios.put(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/set-trailer?ref=${this.save.attached_trailer}`)
        .then(() => this.showToast('Success', 'Trailer attached', 'success'))
        .catch(err => {
          this.showToast('Uhoh…', 'Could not attach trailer', 'danger')
          console.error(err)
        })
    },

    companyNameFromReference(ref) {
      return `${this.companies[ref].city}, ${this.companies[ref].name}`
    },

    createJob() {
      if (!this.companies[this.newJob.origin_reference]) {
        this.showToast('Uhm…', 'Source Company does not exist', 'danger')
        return
      }

      if (!this.companies[this.newJob.target_reference]) {
        this.showToast('Uhm…', 'Target Company does not exist', 'danger')
        return
      }

      if (!this.cargo[this.newJob.cargo_reference]) {
        this.showToast('Uhm…', 'Cargo does not exist', 'danger')
        return
      }

      this.newJob.weight = parseInt(this.newJob.weight)
      if (this.newJob.weight > 200) {
        this.showToast('Uhm…', 'You want to pull more than 200 Tons of cargo?', 'danger')
        return
      }

      this.plannedRoute.push(this.newJob)
      this.newJob = { weight: 10 } // Reset job
    },

    createRoute() {
      this.showSaveModal = true
      return axios.post(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/jobs`, this.plannedRoute)
        .then(() => {
          this.showToast('Success', 'Route created', 'success')
          this.plannedRoute = []
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not add route', 'danger')
          console.error(err)
        })
    },

    executeRepair(fixType) {
      this.showSaveModal = true
      return axios.put(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/fix?type=${fixType}`)
        .then(() => this.showToast('Success', 'Repair executed', 'success'))
        .catch(err => {
          this.showToast('Uhoh…', 'Could not repair', 'danger')
          console.error(err)
        })
    },

    loadCargo() {
      return axios.get(`/api/gameinfo/cargo`)
        .then(resp => {
          this.cargo = resp.data
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not load cargo defintion', 'danger'
            console.error(err)
          })
    },

    loadCompanies() {
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/companies`)
        .then(resp => {
          this.companies = resp.data
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not load company defintion', 'danger')
          console.error(err)
        })
    },

    loadJobs() {
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/jobs`)
        .then(resp => {
          this.jobs = resp.data
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not load jobs', 'danger')
          console.error(err)
        })
    },

    loadNewestSave() {
      this.selectSave(this.sortedSaves[0].id, null)
    },

    loadProfiles() {
      return axios.get('/api/profiles')
        .then(resp => {
          this.profiles = resp.data
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not load profiles', 'danger')
          console.error(err)
        })
    },

    loadSave() {
      // Load companies in background
      this.loadCompanies()

      this.saveLoading = true
      return axios.get(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}`)
        .then(resp => {
          this.save = resp.data
          this.saveLoading = false
        })
        .catch(err => {
          this.showToast('Uhoh…', 'Could not load save-game', 'danger')
          console.error(err)
        })
    },

    loadSaves() {
      if (this.socket) {
        // Dispose old socket
        this.socket.close()
        this.socket = null
      }

      const loc = window.location
      const socketBase = `${loc.protocol === 'https:' ? 'wss:' : 'ws:'}//${loc.host}`
      this.socket = new WebSocket(`${socketBase}/api/profiles/${this.selectedProfile}/saves?subscribe=true`)
      this.socket.onclose = () => window.setTimeout(this.loadSaves, 1000) // Restart socket
      this.socket.onmessage = evt => {
        this.saves = JSON.parse(evt.data)
        this.showSaveModal = false

        if (this.autoLoad) {
          this.loadNewestSave()
        }
      }
    },

    removeJob(idx) {
      if (idx < 0 || idx > this.plannedRoute.length - 1) {
        return
      }

      let newRoute = []
      for (const i in this.plannedRoute) {
        if (parseInt(i) !== idx) {
          newRoute.push(this.plannedRoute[i])
        }
      }

      this.plannedRoute = newRoute
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

    selectSave(saveID) {
      if (this.selectedSave === saveID) {
        this.loadSave()
      } else {
        this.selectedSave = saveID
      }
    },

    setEconomy(param, value) {
      this.showSaveModal = true
      return axios.put(`/api/profiles/${this.selectedProfile}/saves/${this.selectedSave}/economy?${param}=${value}`)
        .then(() => this.showToast('Success', 'Economy updated', 'success'))
        .catch(err => {
          this.showToast('Uhoh…', 'Could not update economy', 'danger')
          console.error(err)
        })
    },

    showToast(title, text, variant = 'info') {
      this.$bvToast.toast(text, {
        'appendToast': true,
        'autoHideDelay': 2500,
        'is-status': true,
        'solid': true,
        title,
        variant,
      })
    },

    validInt(v) {
      return v >= 0 && v <= 2147483647
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
