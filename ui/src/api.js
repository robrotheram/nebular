import axios from 'axios'

class API {
    constructor() {
        this.baseurl = ""
        if (process.env.NODE_ENV === 'development') {
            this.baseurl = "http://localhost:8081"
        }
    }
    getAll(){
        return axios.get(this.baseurl+"/search").then(response => {
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }
    search(term){
        return axios.get(this.baseurl+"/search/"+term).then(response => {
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }
    create(repo){
        return axios.post(this.baseurl+"/roles", repo).then(response => {            
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }
    delete(id){
        return axios.delete(this.baseurl+"/roles/"+id).then(response => {            
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }
    update(id){
        return axios.get(this.baseurl+"/roles/update/"+id).then(response => {            
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }
    user(){
        return axios.get(this.baseurl+"/user").then(response => {
            // returning the data here allows the caller to get it through another .then(...)
            return (response.data)
        })
    }

}
const api = new API()
export {api}