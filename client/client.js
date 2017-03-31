var apiClient = (function() {
  return {
    Config: function(apiAddr, authToken)
    {
      return {
        apiAddr: apiAddr,
      }
    },  
      
    getClinics: function(config, callback)
    {
      var elems = [config.apiAddr, 'clinics']
      this.sendXMLHttpRequest(config, 'GET', elems.join('/'), '', callback)
    },

    addClinic: function(config, data, callback)
    {
      var elems = [config.apiAddr, 'clinics']
      this.sendXMLHttpRequest(config, 'POST', elems.join('/'), JSON.stringify(data), callback)
    },

    getWaitTime: function(config, clinicId, callback)
    {
      var elems = [config.apiAddr, 'clinics', clinicId, 'wait_time']
      this.sendXMLHttpRequest(config, 'GET', elems.join('/'), '', callback)
    },
    
    getPatients: function(config, clinicId, callback)
    {
      var elems = [config.apiAddr, 'clinics', clinicId, 'patients']
      this.sendXMLHttpRequest(config, 'GET', elems.join('/'), '', callback)
    },
    
    addPatient: function(config, clinicId, data, callback)
    {
      var elems = [config.apiAddr, 'clinics', clinicId, 'patients']
      this.sendXMLHttpRequest(config, 'POST', elems.join('/'), JSON.stringify(data), callback)
    },

	removePatient : function(config, clinicId, patientNo, callback)
    {
      var elems = [config.apiAddr, 'clinics', clinicId, 'patients', patientNo]
      this.sendXMLHttpRequest(config, 'DELETE', elems.join('/'), '', callback)
    },

    sendXMLHttpRequest: function(config, method, route, data, callback) {  
        if (typeof callback != 'function') {
          callback = function(r) {
            console.log(r)
          }
        }
        
        var xhr = new XMLHttpRequest()
        
        // asyhnchronous request
        xhr.open(method, route, true)
        
        xhr.onreadystatechange = function() {
          if(xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
              callback(JSON.parse(xhr.responseText))
            } else {
              callback({ error: "unexpected status code", data: xhr })
            }
          }
        }
        
        // send the data
        xhr.send(data)
    },
  } // return
})() 