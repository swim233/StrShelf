<script setup lang="ts">
import { computed, onMounted, ref, TransitionGroup } from 'vue'
import './components/views/Icons.vue'
import Icons from './components/views/Icons.vue'
import Cookies from 'js-cookie'

interface StrShelfResponse<T> {
  code: number
  data: T
  msg: string
}
interface ShelfItem {
  id: number
  title: string
  link: string
  comment: string
  gmt_created: number
  gmt_modified: number
  gmt_deleted: number
  deleted: boolean
}
interface DisplayData {
  saveDatas: ShelfItem[]
  date: Date
}

interface PostResult {
  code: number
  result: any
}
type noticeType = 'Info' | 'Success' | 'Warning' | 'Error'

interface Notice {
  id?: number
  type: noticeType
  delayTime: number
  mainMessage: string
  subMessage?: string
}

interface Token {
  token: string
}

interface UserInfo {
  username: string
  password: string
}

let noticeId = 0

const baseUrl =
  process.env['NODE_ENV'] === 'development' ? 'http://127.0.0.1:1111' : ''

const saveDatas = ref<ShelfItem[]>([])
onMounted(() => {
  notify({
    type: 'Info',
    delayTime: 3000,
    mainMessage: '正在获取列表',
  })
  fetchData()
  loginDisplayText.value = computeDisplayText()
})

const fetchData = async () => {
  let data = await fetch(baseUrl + '/v1/item.get', {
    method: 'POST',
  })
  let readItems = (await data.json()) as StrShelfResponse<ShelfItem[]>
  saveDatas.value = readItems.data
  if (data.ok) {
    notify({
      type: 'Success',
      delayTime: 5000,
      mainMessage: '获取成功',
    })
  } else {
    notify({
      type: 'Error',
      delayTime: 5000,
      mainMessage: '获取失败',
      subMessage: '无法从: ' + data.url + ' 获取列表 请检查网络',
    })
  }
}
//mocker
const mockDate = () => {
  let id: number = 0
  let testDate: ShelfItem = {
    id: id,
    title: '标题',
    link: 'http://example.com',
    comment: '评论',
    gmt_created: new Date().getTime(),
    gmt_deleted: new Date().getTime(),
    gmt_modified: new Date().getTime(),
    deleted: false,
  }
  id++
  saveDatas.value.push(testDate)
  console.log(displayDatas.value)
}

function mapDate(date: Date): string {
  return date.getDate() + '_' + date.getMonth()
}

const displayDatas = computed<DisplayData[]>(() => {
  return Object.values(
    saveDatas.value.reduce<Record<string, DisplayData>>((result, saveData) => {
      const key = mapDate(new Date(saveData.gmt_created))
      if (!result[key])
        result[key] = {
          saveDatas: [saveData],
          date: new Date(saveData.gmt_created),
        }
      else result[key].saveDatas.push(saveData)
      return result
    }, {}),
  )
})

const notices = ref<Notice[]>([])

const notify = (no: Notice) => {
  no.id = noticeId++
  notices.value.push(no)
  setTimeout(() => {
    notices.value.splice(notices.value.indexOf(no), 1)
  }, no.delayTime)
}

const postNewData = async () => {
  let postData: ShelfItem = {
    id: 0,
    title: newData.value.title,
    link: newData.value.link,
    comment: newData.value.comment,
    gmt_created: 0,
    gmt_modified: 0,
    gmt_deleted: 0,
    deleted: false,
  }
  postIsActive.value = false
  dialogDisplayDelay.value = false

  try {
    let response = await fetch(baseUrl + '/v1/item.post', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + Cookies.get('token'),
      },
      body: JSON.stringify(postData),
    })

    let r = (await response.json()) as PostResult
    if (r.result) {
      notify({
        type: 'Success',
        delayTime: 3000,
        mainMessage: '创建成功',
        subMessage: '项目创建成功',
      })
      fetchData()
    } else if (r.code == 401) {
      notify({
        type: 'Error',
        delayTime: 3000,
        mainMessage: '身份验证失败',
        subMessage: '请重新登陆',
      })
      loginStatus.value = false
    } else {
      notify({
        type: 'Error',
        delayTime: 5000,
        mainMessage: '创建失败',
        subMessage: '项目创建失败 请检查网络',
      })
    }
  } catch (error) {
    notify({
      type: 'Error',
      delayTime: 5000,
      mainMessage: '创建失败',
      subMessage: '项目创建失败 请检查网络',
    })
  }
}

const newData = ref<ShelfItem>({
  id: 0,
  title: '',
  link: '',
  comment: '',
  gmt_created: 0,
  gmt_modified: 0,
  gmt_deleted: 0,
  deleted: false,
})

const searchString = ref<string>('')

//reset
const reset = () => {
  saveDatas.value = []
}

const postIsActive = ref<boolean>(false)
const loginIsActive = ref<boolean>(false)
let timeout: ReturnType<typeof setTimeout> | null = null

const post = () => {
  postIsActive.value = true

  if (timeout) {
    clearTimeout(timeout)
  }

  timeout = setTimeout(() => {
    dialogDisplayDelay.value = true
  }, 300)
}
const loginOverlay = ref<boolean>(false)
const loginButton = () => {
  loginOverlay.value = true
  loginIsActive.value = true
}
const userInfo = ref<UserInfo>({
  username: '',
  password: '',
})

const loginSubmit = async () => {
  let result = await login(userInfo.value.username, userInfo.value.password)
  if (result) {
    notify({
      type: 'Success',
      delayTime: 3000,
      mainMessage: '登陆成功',
    })
    loginDisplayText.value = userInfo.value.username
  } else {
    notify({
      type: 'Error',
      delayTime: 3000,
      mainMessage: '登陆失败',
      subMessage: '请检查用户名和密码',
    })
  }
  loginIsActive.value = false
  userInfo.value.password = ''
  userInfo.value.username = ''
  //TODO:更优雅的置空方式
}
const cancelDialog = () => {
  postIsActive.value = false
  loginIsActive.value = false
  dialogDisplayDelay.value = false

  if (timeout) {
    clearTimeout(timeout)
    timeout = null
  }
}
const activeCss = computed<string>(() => {
  return postIsActive.value ? 'post-dialog-active' : 'post-dialog-positive'
})
const loginActiveCss = computed<string>(() => {
  return loginIsActive.value ? 'login-dialog-active' : 'login-dialog-positive'
})

const dialogDisplayDelay = ref<boolean>(false)

const login = async (account: string, password: string): Promise<boolean> => {
  let token = Cookies.get('token')
  console.log('token from cookies: ' + token)
  let isTokenValid = await verifyJWT(token)
  if (isTokenValid) {
    return true
  }
  let response = await fetch(baseUrl + '/v1/user.login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ account, password }),
  })

  if (response.ok) {
    let data = (await response.json()) as Token
    if (await verifyJWT(data.token)) {
      console.log('verify success before set token')
      Cookies.set('token', data.token)
      Cookies.set('username', account)
      return true
    } else {
      console.log('fail to verifyJWT!: ' + data.token)
      console.log(response.status)
      return false
    }
  }
  return false
}

const verifyJWT = async (token: string | undefined): Promise<boolean> => {
  if (typeof token === undefined) {
    return false
  }
  try {
    let response = await fetch(baseUrl + '/v1/user.verify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token }),
    })
    if (response.ok) {
      return true
    }
  } catch (error) {
    console.log(error)
    return false
  }
  return false
}
const loginDisplayText = ref<string>('登陆')

const computeDisplayText = (): string => {
  let username = Cookies.get('username')
  if (username != null) {
    return username
  } else {
    return '登陆'
  }
}

const loginStatus = ref<boolean>(false)
</script>

<template>
  <div class="root">
    <Transition name="masker">
      <div
        class="masker"
        @click="cancelDialog"
        v-if="postIsActive || loginIsActive"
      ></div>
    </Transition>
    <!-- <div class="debug-dialog">debug test</div> -->
    <div class="login-wrapper" v-if="!loginStatus">
      <button class="login-button" @click="loginButton()">
        {{ loginDisplayText }}
      </button>
    </div>
    <div class="login-dialog" :class="loginActiveCss">
      <div class="login-dialog-banner">登陆账号</div>
      <div class="login-dialog-username">
        <input
          class="login-dialog-username-input login-dialog-input"
          placeholder="账号"
          v-model="userInfo.username"
        />
      </div>
      <div class="login-dialog-password">
        <input
          class="login-dialog-password-input login-dialog-input"
          placeholder="密码"
          v-model="userInfo.password"
        />
      </div>
      <div class="login-dialog-submit">
        <button class="login-dialog-submit-button" @click="loginSubmit()">
          登陆
        </button>
      </div>
    </div>
    <div class="post-dialog" :class="activeCss">
      <Transition name="post-dialog-transition">
        <div v-if="dialogDisplayDelay" class="post-dialog-wrapper">
          <div class="post-dialog-banner">新建项目</div>
          <div class="post-dialog-content">
            <span class="post-dialog-content-name-wrapper">
              标题
              <input
                class="post-dialog-content-input post-dialog-title-input"
                placeholder=""
                v-model="newData.title"
              />
              <br />
            </span>
            <span class="post-dialog-content-link-wrapper">
              链接
              <input
                class="post-dialog-content-input post-dialog-link-input"
                v-model="newData.link"
              />
              <br />
            </span>
            <span class="post-dialog-content-comment-wrapper"
              >评论
              <input
                placeholder=""
                class="post-dialog-content-input"
                v-model="newData.comment"
              />
            </span>
          </div>
          <div class="post-dialog-button-wrapper">
            <button class="post-dialog-button" @click="postNewData()">
              提交
            </button>
          </div>
        </div>
      </Transition>
    </div>

    <div class="search-wrapper">
      <input class="search" id="search" placeholder="" v-model="searchString" />
      <label class="search-label" for="search">搜索</label>
      <!-- <p>test:{{ searchString }}</p> -->
    </div>
    <div class="title-wrapper">
      <h1 class="title">StrShelf</h1>
    </div>
    <div class="content-wrapper">
      <TransitionGroup name="content-transition">
        <div
          class="content"
          v-for="displayData in displayDatas"
          :key="mapDate(displayData.date)"
        >
          <div class="date">
            <div class="year">
              {{ displayData.date.getFullYear()
              }}<span class="year-display">年</span>
            </div>
            <div class="day">{{ displayData.date.getDate() }}</div>
            <div class="month">{{ displayData.date.getMonth() + 1 }}月</div>
          </div>
          <div class="edit">
            <button @click="">
              <Icons class="edit-icon" :icon="'Edit'" />
            </button>
          </div>
          <div class="delete">
            <Icons class="delete-icon" :icon="'delete'" />
          </div>
          <div>
            <div
              class="content-inner"
              v-for="data in displayData.saveDatas"
              :key="data.id"
            >
              <div class="content-inner-title">{{ data.title }}</div>
              <a class="content-inner-link" :href="data.link">{{
                data.link
              }}</a>
              <div class="content-inner-comment">{{ data.comment }}</div>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>
    <div class="control-wrapper">
      <div class="control">
        <button @click="mockDate" class="mock">mock</button>
      </div>
      <div class="control">
        <button @click="reset" class="reset">reset</button>
      </div>
      <div class="control">
        <button @click="fetchData" class="fetch">fetch</button>
      </div>
      <div class="control">
        <button @click="post" class="post" v-if="!postIsActive">
          <span class="post-button">+</span>
        </button>
      </div>
    </div>
    <div class="notice-dialog-wrapper">
      <TransitionGroup name="notice-dialog-instance-transition">
        <div
          v-for="notice in notices"
          :class="`notice-dialog-instance ${notice.type}`"
          :key="notice.id"
        >
          <div :class="`notice-dialog-main ${notice.type}`">
            <Icons :icon="notice.type" /> {{ notice.mainMessage }}
          </div>
          <div :class="`notice-dialog-sub ${notice.type}`">
            {{ notice.subMessage }}
          </div>
        </div>
      </TransitionGroup>
    </div>
  </div>
</template>

<style>
:root {
  --background-alpha-color: rgba(8, 22, 61, 0.635);
  --background-color: #14181d;
  --control-animation-duration: 0.4s;
  --animation-color: rgba(127, 127, 127, 0.598);
  --color: #ced4d9;
  --color-button-border: rgb(46, 65, 122);
  --color-green-border: #03482c;
  --color-green-bg: #012012;
  --color-red-border: #481403;
  --color-red-bg: #240401;
  --color-blue-border: rgb(32, 48, 114);
  --color-blue-bg: rgb(19, 28, 73);
  background-color: var(--background-color);
  color: var(--color);
}

.mock {
  position: fixed;
  font-size: 24px;
  bottom: 0;
  margin: 24px;
  min-width: 90px;
  border: 5px solid rgb(40, 40, 141);
  border-radius: 6px;
  background-color: rgba(112, 112, 211, 0.329);
  padding: 6px;
  transition: all var(--control-animation-duration);
  transition-delay: 0.05s;
  backdrop-filter: blur(4px);
}
.reset {
  position: fixed;
  font-size: 24px;
  bottom: 0;
  left: 110px;
  min-width: 90px;
  margin: 24px;
  border: 5px solid rgb(213, 21, 82);
  border-radius: 6px;
  background-color: rgba(206, 27, 102, 0.521);
  padding: 6px;
  transition: all var(--control-animation-duration);
  backdrop-filter: blur(4px);
}
.fetch {
  position: fixed;
  font-size: 24px;
  bottom: 0;
  left: 220px;
  min-width: 90px;
  margin: 24px;
  border: 5px solid rgb(21, 213, 43);
  border-radius: 6px;
  background-color: rgba(23, 177, 20, 0.521);
  padding: 6px;
  transition: all var(--control-animation-duration);
  backdrop-filter: blur(4px);
}

.mock:hover {
  font-size: 30px;
  box-shadow: 0 0 36px rgb(40, 40, 141);
}
.reset:hover {
  font-size: 30px;
  box-shadow: 0 0 36px rgb(213, 21, 82);
}
.fetch:hover {
  font-size: 30px;
  box-shadow: 0 0 36px rgb(21, 213, 43);
}
.control > * {
  font-weight: bold;
}
.title-wrapper {
  display: flex;
}
.title:hover {
  text-shadow: 0 0 24px var(--animation-color);
}
.title {
  font-size: 50px;
  font-weight: 500;
  margin: 36px;
  transition-property: all;
  transition-duration: 0.5s;
  transition-delay: 0.1s;
}

.content-wrapper {
  display: flex;
  flex-direction: column;
}

.date {
  display: flex;
}
.year {
  display: flex;
  margin-left: 24px;
  margin-right: 0;
  font-size: 24px;
  font-weight: 500;
  transition: all 0.5s;
  transition-delay: 0.1s;
  opacity: 0;
}
.year:hover {
  text-shadow: 0 0 24px var(--animation-color);
}
.year-display {
  font-size: small;
  margin-top: 3px;
}
.content:hover .year {
  opacity: 1;
}
.month {
  width: 60px;
  font-size: 24px;
  font-weight: 500;
  min-width: 90px;
  margin-top: 6px;
  transition: all 0.5s;
  transition-delay: 0.1s;
}
.month:hover {
  text-shadow: 0 0 24px var(--animation-color);
}
.day {
  display: flex;
  justify-content: flex-end;
  width: 60px;
  font-size: 48px;
  font-weight: 500;
  margin: 0 30px 0 0;
  transition: all 0.5s;
  transition-delay: 0.1s;
}
.day:hover {
  text-shadow: 0 0 24px var(--animation-color);
}

.edit,
.delete {
  transition: all 0.5s;
  transition-delay: 0.1s;
  margin-top: 6px;
  opacity: 0;
}
.edit {
  margin-right: 6px;
}
.edit-icon {
  color: rgb(83, 124, 221);
}
.delete-icon {
  color: #8d2832;
}
.content:hover .edit,
.content:hover .delete {
  opacity: 1;
}

.content {
  display: flex;
}
.content-inner {
  margin-left: 30px;
  margin-top: 6px;
}
.content-inner-comment {
  display: flex;
  margin-bottom: 24px;
  font-size: 16px;
  transition: all 0.5s;
  transition-delay: 0.1s;
}
.content-inner-comment:hover {
  text-shadow: 0 0 24px var(--animation-color);
}
.content-inner-link {
  display: flex;
  margin-bottom: 12px;
  font-size: 16px;
  transition: all 0.5s;
  transition-delay: 0.1s;
  color: var(--color);
}
.content-inner-link:hover {
  text-shadow: 0 0 24px var(--animation-color);
}
.content-inner-title {
  display: flex;
  margin-bottom: 12px;
  font-size: 24px;
  transition: all 0.5s;
  transition-delay: 0.1s;
}
.content-inner-title:hover {
  text-shadow: 0 0 24px var(--animation-color);
}

.search-wrapper {
  background: none;
  margin: 36px 0 0 36px;
}

.search {
  color: var(--color);
  position: relative;
  padding: 12px;
  border: 3px solid var(--color);
  border-radius: 6px;
  outline: none;
  transition: box-shadow 0.4s 0.1s;
  font-size: large;
  font-weight: bold;
  border-color: #90abea;
}

.search::placeholder {
  color: var(--color);
}
.search-label {
  position: absolute;
  top: 50px;
  left: 45px;
  cursor: text;
  user-select: none;
  font-weight: bold;
  padding: 0 0.2em;
  font-size: larger;
  transform-origin: top left;
  transition: all 0.4s 0.1s;
}
.search:hover ~ .search-label,
.search:focus ~ .search-label,
.search:not(:placeholder-shown) ~ .search-label {
  font-size: small;
  transform: translateY(-140%) translateX(-5%);
  background-color: var(--background-color);
}

.content-transition-move,
.content-transition-enter-active,
.content-transition-leave-active {
  transition: all 0.3s ease;
}
.content-transition-leave-active {
  position: absolute;
}
.content-transition-enter-from {
  opacity: 0;
  transform: translateY(30px);
}
.content-transition-leave-to {
  opacity: 0;
}
.masker-enter-active,
.masker-leave-active {
  transition: opacity 0.3s;
}

.masker-enter-from,
.masker-leave-to {
  opacity: 0 !important;
}
/*
.masker-enter-to,
.masker-leave-from {
  opacity: 0.8;
} */

.masker {
  transition: all 0.3s;
  backdrop-filter: blur(4px);
  position: fixed;
  background-color: rgba(0, 0, 0, 0.6);
  width: 100vw;
  height: 100vh;
  left: 0;
  top: 0;
  z-index: 1;
}

.post,
.post-dialog {
  color: white;
  position: fixed;
  border: 5px solid var(--color-button-border);
  border-radius: 8px;
  background-color: #516a89;
  width: 60px;
  height: 60px;
  bottom: 10px;
  right: 10px;
  z-index: 2;
  transition: all 0.3s;
  transform-origin: bottom right;
  /* width 0.3s,
    height 0.3s,
    transform 0.3s 0.1s !important; */
}
.post {
  font-size: 200%;
  opacity: 0;
}
/* .post-dialog-positive {
  opacity: 0;
} */

.post-dialog-active {
  /* transform: scale(500%, 500%) translate(-200%, -100%); */
  display: flex;
  width: 70vw;
  height: 60vh;
  transform: translate(-20%, -20%);
  opacity: 1;
  color: #14181d;
  background-color: #14181d;
}

.post-dialog-wrapper {
  color: #ced4d9;
  width: 100%;
  height: 100%;
  font-size: large;
  font-weight: bold;
}

.post-dialog-banner {
  font-size: 30px;
  display: flex;
  justify-content: center;
  font-weight: bold;
  margin: 20px 20px;
}

.post-dialog-content * {
  margin: 20px 0 0 20px;
  font-weight: 800;
  font-size: 24px;
}
.post-dialog-content-input {
  color: var(--color);
  position: relative;
  padding: 12px;
  border: 3px solid var(--color);
  border-radius: 6px;
  outline: none;
  transition: box-shadow 0.4s 0.1s;
  font-size: large;
  font-weight: bold;
}

.post-dialog-title-input:required:invalid,
.post-dialog-link-input:required:invalid {
  border-color: var(--color-red-border);
}
.post-dialog-content-input::placeholder {
  color: var(--color);
}

.post-dialog-button-wrapper {
  position: absolute;
  height: 60px;
  width: 90px;
  bottom: 30px;
  right: 30px;
  border: 4px solid rgb(22, 30, 62);
  border-radius: 8px;
}
.post-dialog-button {
  border: 1px solid rgb(22, 30, 62);
  border-radius: 8px;
  display: flex;
  height: 100%;
  width: 100%;
  justify-content: center;
  align-items: center;
  color: white;
  font-weight: 800;
  font-size: 24px;
  background-color: #516a89;
}

.notice-dialog-wrapper {
  z-index: 3;
  box-sizing: content-box;
  top: 30px;
  right: 30px;
  pointer-events: none;
  /* min-width: 24vw; */
  /* min-height: 12vh; */
  position: fixed;
}

.notice-dialog-instance {
  padding-right: 20px;
  display: flex;
  min-height: 12vh;
  min-width: 20vw;
  flex-direction: column;
  margin-bottom: 12px;
}
.notice-dialog-instance.Success {
  border: 5px solid var(--color-green-border);
  border-radius: 8px;
  background-color: var(--color-green-bg);
}
.notice-dialog-instance.Error {
  border: 5px solid var(--color-red-border);
  border-radius: 8px;
  background-color: var(--color-red-bg);
}
.notice-dialog-instance.Info {
  border: 5px solid var(--color-blue-border);
  border-radius: 8px;
  background-color: var(--color-blue-bg);
}

.notice-dialog-main {
  position: relative;
  margin-top: 12px;
  left: 12px;
  font-weight: 700;
  font-size: 24px;
}

.notice-dialog-sub {
  color: #9e9e9e;
  position: relative;
  font-size: 14px;
  margin-top: 12px;
  margin-bottom: 24px;
  width: 100%;
  font-weight: 700;
  text-indent: 1em;
}

.notice-dialog-instance-transition-move,
.notice-dialog-instance-transition-enter-active,
.notice-dialog-instance-transition-leave-active {
  transition: all 0.3s ease;
}
.notice-dialog-instance-transition-leave-active {
  position: absolute;
}
.notice-dialog-instance-transition-enter-from {
  opacity: 0;
}
.notice-dialog-instance-transition-leave-to {
  opacity: 0;
}

.post-dialog-transition-move,
.post-dialog-transition-enter-active {
  transition: all 0.3s ease;
}
.post-dialog-transition-leave-active {
  position: absolute;
}
.post-dialog-transition-enter-from {
  opacity: 0;
}
.post-dialog-transition-leave-to {
  opacity: 0;
}
.login-wrapper,
.login-button {
  position: fixed;
  top: 30px;
  right: 30px;
  font-size: 18px;
  font-weight: bolder;
  color: aliceblue;
}

.login-dialog-positive {
  opacity: 0;
  pointer-events: none;
}

.login-dialog {
  display: flex;
  position: fixed;
  border: 3px solid var(--color-button-border);
  border-radius: 6px;
  height: 35vh;
  width: 30vw;
  top: 20vh;
  left: 35vw;
  background-color: #14181d;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  z-index: 2;
  background-color: var(--background-alpha-color);
}

.login-dialog-banner {
  font-size: 18px;
  display: flex;
  justify-content: center;
  font-weight: bold;
  margin: 18px 20px;
}
.login-dialog-input {
  color: var(--color);
  position: relative;
  padding: 6px;
  border: 3px solid var(--color);
  border-radius: 6px;
  outline: none;
  transition: box-shadow 0.4s 0.1s;
  font-size: large;
  margin: 3px;
}

.login-dialog-submit-button {
  margin: 12px;
  border: 3px solid var(--color-green-border);
  border-radius: 6px;
  padding: 6px;
  color: aliceblue;
  font-size: large;
  background-color: var(--color-green-bg);
}
.debug-dialog {
  z-index: 10;
  font-size: 30px;
  width: fit-content;

  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

*::selection {
  color: var(--color);
}
</style>
