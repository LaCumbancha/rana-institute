import http from 'k6/http';
import { check, sleep, group } from "k6"
import { Trend, Rate, Counter } from 'k6/metrics'
import papaparse from 'https://jslib.k6.io/papaparse/5.1.1/index.js'

// Setting test URLs
const DEV_URL = "http://localhost:8080/"
const PROD_URL = "https://site-dot-taller3-rana.uc.r.appspot.com/"

let URL = DEV_URL
if (__ENV.ENVIRONMENT !== undefined) {
	switch (__ENV.ENVIRONMENT.toUpperCase()) {
		case 'PROD':
			URL = PROD_URL
			break
		default:
			URL = DEV_URL
	}
}

// Setting pages
let currentPage = "/"
const IDX_PAGE = "/"
const HOME_PAGE = "/home"
const JOBS_PAGE = "/jobs"
const ABOUT_PAGE = "/about"
const LEGAL_PAGE = "/about/legal"

// Setting cached resources
const FAVICON_URL = `${URL}/favicon.ico`
const HOME_IMG_URL = `${URL}/static/mister-chispa.webp`
const JOBS_IMG_URL = `${URL}/static/solitaire-bob.png`
const ABOUT_IMG_URL = `${URL}/static/happy-student.jpeg`
const LEGAL_IMG1_URL = `${URL}/static/bell.gif`
const LEGAL_IMG2_URL = `${URL}/static/clock.gif`
const LEGAL_IMG3_URL = `${URL}/static/mouth.gif`
const LEGAL_IMG4_URL = `${URL}/static/toaster.gif`
const LEGAL_IMG5_URL = `${URL}/static/worm.gif`
const LEGAL_IMG6_URL = `${URL}/static/dancing-jesus.gif`

let userCachedResources = { IDX_PAGE: false, HOME_PAGE: false, JOBS_PAGE: false, ABOUT_PAGE: false, LEGAL_PAGE: false }

// Setting metrics
const errorRate = new Rate('ErrorRate')
const successRate = new Rate('SuccessRate')
let waitingTime = new Trend('WaitingTime')
let statusCodes = new Counter('StatusCodes')
let failingRequests = new Counter('FailingRequests')

// Loading test stages
const csvData = papaparse.parse(open(__ENV.SCENARIO), { header: true }).data
let csvStages = []

csvData.forEach(function(csvRow) {
	if (csvRow.duration != "") {
		csvStages.push({ duration: csvRow.duration, target: parseInt(csvRow.target) })
	}
})

// Setting K6 options
export let options = {
	discardResponseBodies: false,
	scenarios: {
		stressTest: {
			executor: 'ramping-vus',
			exec: 'stressTest',
			startVUs: 0,
			stages: csvStages,
			gracefulRampDown: '10s',
		},
	}
}

// Updating metrics
export function metrics(status, time) {
	waitingTime.add(time)

	if (status < 200 || status >= 299) {
		failingRequests.add(1)
	}

	statusCodes.add(1, { tag: status })
	errorRate.add(status < 200 || status >= 300)
	successRate.add(status >= 200 && status < 300)
    check(status, { 'Status was 2XX': (code) => code >= 200 && code < 300 } )
}

// Test run
export function stressTest() {
	group('/site', function(){
		currentPage = nextPage(currentPage)
		let response = http.get(URL + currentPage)
		let responseTime = response.timings.waiting

		// Retrieving favicon only once.
		if (!userCachedResources[IDX_PAGE]) {
			userCachedResources[IDX_PAGE] = true
			let imgResponse = http.get(FAVICON_URL)
			responseTime += imgResponse.timings.waiting
		}

		// Retrieving static content for HOME page (only once).
		if (currentPage == HOME_PAGE) {
			if (!userCachedResources[HOME_PAGE]) {
				userCachedResources[HOME_PAGE] = true
				let imgResponse = http.get(HOME_IMG_URL)
				responseTime += imgResponse.timings.waiting
			}
		}

		// Retrieving static content for JOBS page (only once).
		if (currentPage == JOBS_PAGE) {
			if (!userCachedResources[JOBS_PAGE]) {
				userCachedResources[JOBS_PAGE] = true
				let imgResponse = http.get(JOBS_IMG_URL)
				responseTime += imgResponse.timings.waiting
			}
		}

		// Retrieving static content for ABOUT page (only once).
		if (currentPage == ABOUT_PAGE) {
			if (!userCachedResources[ABOUT_PAGE]) {
				userCachedResources[ABOUT_PAGE] = true
				let imgResponse = http.get(ABOUT_IMG_URL)
				responseTime += imgResponse.timings.waiting
			}
		}

		// Retrieving static content for LEGAL page (only once).
		if (currentPage == LEGAL_PAGE) {
			if (!userCachedResources[LEGAL_PAGE]) {
				userCachedResources[LEGAL_PAGE] = true
				let imgResponse1 = http.get(LEGAL_IMG1_URL)
				let imgResponse2 = http.get(LEGAL_IMG2_URL)
				let imgResponse3 = http.get(LEGAL_IMG3_URL)
				let imgResponse4 = http.get(LEGAL_IMG4_URL)
				let imgResponse5 = http.get(LEGAL_IMG5_URL)
				let imgResponse6 = http.get(LEGAL_IMG6_URL)
				responseTime += imgResponse1.timings.waiting + imgResponse2.timings.waiting + 
					imgResponse3.timings.waiting + imgResponse4.timings.waiting + 
					imgResponse5.timings.waiting + imgResponse6.timings.waiting
			}
		}

		metrics(response.status, responseTime)
		sleep(Math.random() * 10)
	})
}

// Auxiliar function for randomly select the next page
function nextPage(page) {
	let choice = Math.random()
	switch (page) {
		case IDX_PAGE:
			if (choice < 0.33) { return HOME_PAGE }
			else if (choice < 0.66) { return JOBS_PAGE }
			else { return ABOUT_PAGE }
		case HOME_PAGE:
			if (choice < 0.5) { return JOBS_PAGE }
			else { return ABOUT_PAGE }
		case JOBS_PAGE:
			if (choice < 0.5) { return HOME_PAGE }
			else { return ABOUT_PAGE }
		case ABOUT_PAGE:
			if (choice < 0.2) { return HOME_PAGE }
			else if (choice < 0.4) { return JOBS_PAGE }
			else { return LEGAL_PAGE }
		case LEGAL_PAGE:
			if (choice < 0.33) { return HOME_PAGE }
			else if (choice < 0.66) { return JOBS_PAGE }
			else { return ABOUT_PAGE }
	}
}
