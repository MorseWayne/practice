// 帖子数据存储
let posts = [];

// 创建新帖
function createPost() {
    const title = document.getElementById('postTitle').value;
    const content = document.getElementById('postContent').value;
    
    const newPost = {
        id: Date.now(),
        title,
        content,
        comments: [],
        votes: 0
    };

    posts.unshift(newPost);
    renderPosts();
}

// 渲染帖子列表
function renderPosts() {
    const container = document.getElementById('postsContainer');
    container.innerHTML = posts.map(post => `
        <div class="post">
            <h3>${post.title}</h3>
            <p>${post.content}</p>
            <div class="post-actions">
                <button onclick="addComment(${post.id})">评论</button>
                <span class="votes">${post.votes} 赞</span>
            </div>
            <div class="comments" id="comments-${post.id}"></div>
        </div>
    `).join('');
}

// 模拟问题数据
const questionsData = [
    {
        id: 1,
        title: "如何在JavaScript中实现数组去重？",
        excerpt: "我有一个包含重复元素的数组，想要得到一个没有重复元素的新数组。尝试了几种方法，但效率不高，有没有更优的解决方案？",
        votes: 24,
        answers: 5,
        views: 128,
        tags: ["JavaScript", "数组"],
        author: "张三",
        time: "2小时前"
    },
    {
        id: 2,
        title: "React Hooks如何正确管理组件状态？",
        excerpt: "在使用useState和useEffect时遇到了一些问题，特别是在处理异步操作时。状态更新后组件没有重新渲染，请问如何正确使用React Hooks管理组件状态？",
        votes: 18,
        answers: 3,
        views: 96,
        tags: ["React", "Hooks", "前端"],
        author: "李四",
        time: "5小时前"
    },
    {
        id: 3,
        title: "Python中列表和元组的区别是什么？",
        excerpt: "刚学习Python，不太清楚列表(list)和元组(tuple)之间的主要区别，以及在什么情况下应该使用哪个。有人能详细解释一下吗？",
        votes: 32,
        answers: 8,
        views: 215,
        tags: ["Python", "基础"],
        author: "王五",
        time: "昨天"
    },
    {
        id: 4,
        title: "如何优化MySQL查询性能？",
        excerpt: "我的数据库查询随着数据量增长变得越来越慢，特别是在关联多个表的时候。除了添加索引之外，还有哪些优化MySQL查询性能的方法？",
        votes: 45,
        answers: 12,
        views: 389,
        tags: ["MySQL", "性能优化", "数据库"],
        author: "赵六",
        time: "3天前"
    }
];

// 页面加载完成后执行
document.addEventListener('DOMContentLoaded', function() {
    // 如果是首页，则渲染问题列表
    if (window.location.pathname.endsWith('index.html') || window.location.pathname === '/') {
        renderQuestionsList();
    }
    
    // 如果是提问页面，则初始化表单提交事件
    if (window.location.pathname.endsWith('ask.html')) {
        initQuestionForm();
    }
});

// 渲染问题列表
function renderQuestionsList() {
    const questionsList = document.getElementById('questionsList');
    if (!questionsList) return;

    questionsList.innerHTML = questionsData.map(question => `
        <div class="question-item">
            <div class="question-header">
                <div class="question-stats">
                    <div class="stat-item">
                        <span class="stat-number">${question.votes}</span>
                        <span class="stat-label">投票</span>
                    </div>
                    <div class="stat-item">
                        <span class="stat-number">${question.answers}</span>
                        <span class="stat-label">回答</span>
                    </div>
                    <div class="stat-item">
                        <span class="stat-number">${question.views}</span>
                        <span class="stat-label">浏览</span>
                    </div>
                </div>
                <div class="question-content">
                    <h3 class="question-title"><a href="#">${question.title}</a></h3>
                    <p class="question-excerpt">${question.excerpt}</p>
                    <div class="question-meta">
                        <div class="question-tags">
                            ${question.tags.map(tag => `<a href="#" class="tag">${tag}</a>`).join('')}
                        </div>
                        <div class="user-info">
                            <div class="user-avatar">${question.author.charAt(0)}</div>
                            <span>${question.author}</span>
                            <span>•</span>
                            <span>${question.time}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

// 初始化提问表单
function initQuestionForm() {
    const questionForm = document.getElementById('question-form');
    if (!questionForm) return;

    questionForm.addEventListener('submit', function(e) {
        e.preventDefault();
        const title = document.getElementById('title').value;
        const content = document.getElementById('content').value;
        const tags = document.getElementById('tags').value;

        // 简单验证
        if (!title.trim() || !content.trim()) {
            alert('标题和内容不能为空');
            return;
        }

        // 这里可以添加表单提交逻辑
        alert('问题提交成功！\n\n标题: ' + title + '\n内容: ' + content.substring(0, 50) + '...');
        window.location.href = 'index.html';
    });
}