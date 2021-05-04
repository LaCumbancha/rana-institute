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

const LOGO_IMG = "LOGO"
const LEGAL_IMG = "LEGAL"
const LOGO_IMG_URL = "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Monsters_University_Logo.svg/1280px-Monsters_University_Logo.svg.png"
const LEGAL_IMG_URL = "https://media.canalnet.tv/2020/01/Screen-Shot-2020-01-23-at-09.27.32.png"
let userCachedResources = { LOGO_IMG: false, LEGAL_IMG: false }

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


let errorRate = new Rate('ErrorRate')
let trend = new Trend('WaitingTime')


export let options = {
	discardResponseBodies: false,
	scenarios: {
		stressTest: {
			executor: 'ramping-vus',
			exec: 'stressTest',
			startVUs: 0,
			stages: [
				{ duration: '10s', target: 2 }
			],
			gracefulRampDown: '10s',
		},
	},
}


export function metrics(status, time) {
	trend.add(time);
	console.log(`${__VU}.${__ITER} - Status: ${status}. Time: ${time}`)
    check(status, { 'Status was 2XX': (code) => code >= 200 && code < 300 })
}

export function stressTest() {
	group('/site', function(){
		currentPage = nextPage(currentPage)
		let response = http.get(URL + currentPage)
		let responseTime = response.timings.waiting

		if (!userCachedResources[LOGO_IMG]) {
			userCachedResources[LOGO_IMG] = true
			let imgResponse = http.get(LOGO_IMG_URL)
			responseTime += imgResponse.timings.waiting
		}

		if (currentPage == LEGAL_PAGE) {
			if (!userCachedResources[LEGAL_IMG]) {
				userCachedResources[LEGAL_IMG] = true
				let imgResponse = http.get(LEGAL_IMG_URL)
				responseTime += imgResponse.timings.waiting
			}
		}

		metrics(response.status, responseTime)
		sleep(Math.random() * 5)
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
