export default {
  header: {
    navigation: {
      whitepapers: '白皮书',
      downloads: '下载区',
      explorer: '区块浏览器',
      blog: '开发日志',
      roadmap: 'Roadmap'
    },
  },
  footer: {
    getStarted: '开始',
    explore: '浏览',
    community: '社区',
    about: '关于',
    wallet_android: '下载Android钱包',
    wallet_mac: '下载MacOS钱包',
    wallet_win: '下载Windows钱包',
    wallet_lin: '下载Linux钱包',
    infographics: '信息图表',
    whitepapers: '白皮书',
    blockchain: '区块浏览器',
    blog: '开发日志',
    twitter: 'Twitter',
    explorer: 'Explorer',
    price: 'CoinPaprika',
    platform: '平台',
    reddit: 'Reddit',
    github: 'GitHub',
    telegram: 'Telegram',
    slack: 'Slack',
    roadmap: 'Roadmap',
    skyMessenger: 'Sky-Messenger',
    cxPlayground: 'CX Playground',
    team: 'Team',
    subscribe: 'Mailing List',
    market: 'Markets',
    bitcoinTalks: 'Bitcointalks ANN',
    instagram: 'Instagram',
    facebook: 'Facebook',
    discord: 'Discord',
  },
  distribution: {
    rate: '价格: 1 {coinType} = {rate} MDL',
    inventory: '现在可以买最多<strong> {coins} MDL</strong>',
    title: 'MDL自动取款机',
    heading: 'MDL自动取款机',
    headingEnded: 'MDL自动取款机暂时关门',
    ended: `<p>请加入<a href="https://t.me/mdl">MDL Telegram</a>，
       也可以看看
      <a href="https://twitter.com/mdlproject">MDL Twitter</a>.`,
    instructions: `<p><a href="https://coinpaprika.com/coin/mdl-mdl">在CoinPaprika看MDL价格</a></p>

<p>参加MDL分发活动:</p>

<ul>
  <li>在下面输入您的MDL地址</li>
  <li>您将收到一个唯一的{coinType}地址用来购买MDL</li>
  <li>将{coinType}发送到您收到的地址上</li>
</ul>

<p>您可以通过输入您的MDL地址并点击下面的"<strong>检查状态</strong>"来核实订单的状态</p>
<p>每次当您点击<strong>获取地址</strong>, 系统会产生一个新的{coinType}地址</p>
<p>一个SKY地址最多只准许兑换{max_bound_addrs}个比特币</p>
    `,
    statusFor: 'SKY地址{mdlAddress}的订单状态',
    enterAddress: '输入MDL地址',
    getAddress: '获取地址',
    checkStatus: '检查状态',
    loading: '加载中...',
    btcAddress: 'SKY地址',
    errors: {
      noSkyAddress: '请输入您的MDL地址',
      coinsSoldOut: 'MDL自动取款机暂时关门',
    },
    statuses: {
      done: '交易 {id}: MDL已经发送并确认(更新于{updated}).',
      waiting_deposit: '交易 {id}: 等待存入(更新于 {updated}).',
      waiting_send: '交易 {id}: 存入已确认; MDL发送在队列中 (更新于 {updated}).',
      waiting_confirm: '交易 {id}: MDL已发送,等待交易确认 (更新于 {updated}).',
    },
  },
};
