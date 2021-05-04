import http from 'k6/http';
import { check, sleep, group } from "k6"
import { Trend, Rate } from 'k6/metrics'

const DEV_URL = "http://localhost:8080/"
const PROD_URL = "https://site-dot-taller3-rana.uc.r.appspot.com/"

const IDX_PAGE = "/"
const HOME_PAGE = "/home"
const JOBS_PAGE = "/jobs"
const ABOUT_PAGE = "/about"
const LEGAL_PAGE = "/about/legal"

let currentPage = "/"

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

const HOME_IMG_URL = `${URL}/static/mister-chispa.webp`
const JOBS_IMG_URL = `${URL}/static/solitaire-bob.png`
const ABOUT_IMG_URL = `${URL}/static/happy-student.jpeg`
const LEGAL_IMG1_URL = `${URL}/static/bell.gif`
const LEGAL_IMG2_URL = `${URL}/static/clock.gif`
const LEGAL_IMG3_URL = `${URL}/static/mouth.gif`
const LEGAL_IMG4_URL = `${URL}/static/toaster.gif`
const LEGAL_IMG5_URL = `${URL}/static/worm.gif`
const LEGAL_IMG6_URL = `${URL}/static/dancing-jesus.gif`

let userCachedResources = { HOME_PAGE: false, JOBS_PAGE: false, ABOUT_PAGE: false, LEGAL_PAGE: false }


let errorRate = new Rate('ErrorRate')
let trend = new Trend('WaitingTime')

let longScenario = {
	stressTest: {
		executor: 'ramping-vus',
		exec: 'stressTest',
		startVUs: 0,
		stages: [
			{ duration: '4m', target: 400 },
			{ duration: '4m', target: 400 },
			{ duration: '2m', target: 0 },
		],
		gracefulRampDown: '10s',
	},
}

let shortScenario = {
	stressTest: {
		executor: 'ramping-vus',
		exec: 'stressTest',
		startVUs: 0,
		stages: [
			{ duration: '30s', target: 5 }
		],
		gracefulRampDown: '10s',
	},
}

let testScenario = shortScenario
if (__ENV.SIZE !== undefined) {
	switch (__ENV.SIZE.toUpperCase()) {
		case 'LONG':
			testScenario = longScenario
			break
		default:
			testScenario = shortScenario
	}
}

export let options = {
	discardResponseBodies: false,
	scenarios: testScenario,
}


export function metrics(status, time) {
	trend.add(time);
    check(status, { 'Status was 2XX': (code) => code >= 200 && code < 300 })
}

export function stressTest() {
	group('/site', function(){
		currentPage = nextPage(currentPage)
		let response = http.get(URL + currentPage)
		let responseTime = response.timings.waiting

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
