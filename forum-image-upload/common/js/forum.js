window.addEventListener('load', function () {
    var elements = document.getElementsByClassName("likes");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", function () {
            likeMe("like", this.closest("tr").getAttribute("sub"), this.closest("tr").getAttribute("main"), this.parentElement)
        });
    }
    var elements = document.getElementsByClassName("dislikes");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", function () {
            likeMe("dislike", this.closest("tr").getAttribute("sub"), this.closest("tr").getAttribute("main"), this.parentElement)
        });
    }
    var elements = document.getElementsByClassName("frozenscreen");
    for (let i = 0; i < elements.length; i++) {
        elements[i].addEventListener("click", function () {
        });
    }

    // for filters
    const usernameElmnt = document.querySelector('p[username]');
    if (usernameElmnt) {
        const username = usernameElmnt.innerHTML
        if (username === "") {
            document.getElementById("FilterCategory").addEventListener("change", filterCategories);
        } else {
            document.getElementById("FilterCategory").addEventListener("change", filterTopics);
            document.getElementById("FilterUserTopics").addEventListener("change", filterTopics);
            document.getElementById("FilterMyLikes").addEventListener("change", filterTopics);
        }
    }
})

function showform(id) {
    const imgInserted = document.getElementById("TopicImg" + id);
    const imgCurrent = document.getElementById("currentImg");
    const CategoryCurrent = document.getElementById("TopicCategory");

    if (imgInserted != null) {
        imgCurrent.classList.remove("hidden");
        const imgName = imgInserted.getAttribute("alt");
        a = imgName.substring(0, imgName.lastIndexOf("-"))
        b = imgName.substring(imgName.lastIndexOf("."))
        document.getElementById("uploadImgName").innerHTML = a + b;
    } else if (imgCurrent != null) {
        document.getElementById("currentImg").classList.add("hidden");
    }
    if (typeof id == "undefined") {
        document.getElementById("TopicTitle").innerHTML = "";
        document.getElementById("TopicText").innerHTML = "";
        document.getElementById("TopicId").value = "";
        if (CategoryCurrent != null) {
            CategoryCurrent.value = "";
        }
    } else {
        document.getElementById("TopicTitle").innerHTML = document.getElementById("Topic" + id).innerHTML;
        document.getElementById("TopicText").innerHTML = document.getElementById("TopicComment" + id).innerHTML;
        document.getElementById("TopicId").value = id;
        if (CategoryCurrent != null) {
            var children = CategoryCurrent.childNodes;
            var categories = document.querySelector('[main="' + id + '"]').querySelector("[categories]").innerHTML + ",";
            children.forEach(function (item) {
                if (categories.includes(item.value + ",") == true) {
                    item.setAttribute("selected", "");
                }
            });
        }
    }
    document.getElementById('NewTopic').classList.remove("hidden");
    const closeFormOnEscapekey = (event) => {
        if (event.key === "Escape") {
            document.getElementById('NewTopic').classList.add("hidden");
            document.removeEventListener("keydown", closeFormOnEscapekey)
        }
    }
    document.addEventListener("keydown", closeFormOnEscapekey)
}

function deleterecord(part, id, main) {
    const text = part == "Topic" ? "Topic" : "Topic content";
    if (confirm(text + " will be deleted!") == true) {
        location.href = "/dtopic?TopicType=" + part + "&Id=" + id + "&MainId=" + main;
    }
}

function validateFile() {
    const fileInput = document.getElementById('uploadImg');
    const fileName = fileInput.files[0].name;
    const fileSize = fileInput.files[0].size;
    const maxSize = 20 * 1024 * 1024;
    const allowedExtensions = /(\.jpeg|\.jpg|\.png|\.gif)$/i;
    document.getElementById("uploadImgName").innerHTML = fileName;

    if (!allowedExtensions.exec(fileName)) {
        document.getElementById('fileErrorMsg').innerHTML = 'Invalid file type. Only JPEG, PNG, and GIF files are allowed.';
        fileInput.value = '';
        return false;
    }
    else if (fileSize > maxSize) {
        document.getElementById('fileErrorMsg').innerHTML = 'File size exceeds 20 MB limit.<div style="font-size:20px; color:red;">File will not be included in topic!</div>';
        fileInput.value = '';
        return false;
    }
    else {
        document.getElementById('fileErrorMsg').innerHTML = '';
        return true;
    }

}

function checkName(ctrl) {
    const xhttp = new XMLHttpRequest();
    xhttp.onload = function () {
        if (this.responseText != "") {
            alert("Username already exists: " + this.responseText + ". Try another one!")
            document.getElementById("login_radio").classList.add("hidden");
        } else {
            document.getElementById("login_radio").classList.remove("hidden");
        }
    }
    xhttp.open("POST", "/checkUser");
    xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xhttp.send("user=" + ctrl.value.trim());
}


function loginShow(ctrl) {
    if (ctrl.value.trim().length == 0) {
        document.getElementById("login_radio").classList.add("hidden");
    } else {
        document.getElementById("login_radio").classList.remove("hidden");
    }
}

function loginType(ctrl) {
    if (ctrl.value == "Password") {
        document.getElementById("password_input").classList.remove("hidden");
        document.getElementById("email_signin").value = "";
        document.getElementById("pass_signin").value = "";
    } else {
        document.getElementById("password_input").classList.add("hidden");
    }
}
function likeMe(type, id, mainid, ctrl) {
    const element = document.getElementById('NewTopic');
    if (typeof (element) != 'undefined' && element != null) {
        const xhttp = new XMLHttpRequest();
        xhttp.onload = function () {
            const values = JSON.parse(this.responseText);
            likes = ctrl.getElementsByClassName('likes');
            dislikes = ctrl.getElementsByClassName('dislikes');
            likes[0].innerHTML = "&#128077; " + values["Likes"];
            likes[0].setAttribute("like", values["MyLike"]);
            dislikes[0].innerHTML = "&#128078; " + values["DisLikes"];
            dislikes[0].setAttribute("dislike", values["MyDisLike"]);
        }
        xhttp.open("POST", "/likeme");
        xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xhttp.send("LikeType=" + type + "&Id=" + id + "&MainId=" + mainid);
    } else {
        alert("You must sign in to like/dislike topics!")
    }
}

function handleCredentialResponse(googleUser) {
    const base64Url = googleUser.credential.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    document.getElementById("email_signin").value = JSON.parse(jsonPayload).email;
    document.getElementById("pass_signin").value = "***************************************";
    document.getElementById("signin").click();
}

function openFilters() {
    if (document.getElementById("filters").classList.contains("hidden")) {
        document.getElementById("filters").classList.remove("hidden")
        document.getElementById("filters_button").value = "Close filters"

    } else {
        document.getElementById("filters").classList.add("hidden")
        document.getElementById("filters_button").value = "Open filters"
    }
}


function filterTopics() {
    const topics = document.getElementsByClassName("topic-group");
    // CATEGORY
    const categorySelectedOptions = Array
        .from(document.getElementById("FilterCategory").options)
        .filter(option => option.selected)
        .map(option => option.value);
    if (categorySelectedOptions.length === 0) { categorySelectedOptions[0] = "All categories" }
    // USERS
    const userSelectedOptions = Array
        .from(document.getElementById("FilterUserTopics").options)
        .filter(option => option.selected)
        .map(option => option.value);
    if (userSelectedOptions.length === 0) { userSelectedOptions[0] = "All users" }
    // LIKES
    const likesSelectedOption = document.getElementById("FilterMyLikes").value;

    for (let i = 0; i < topics.length; i++) {
        const topic = topics[i];
        let categoryShowTopic = false;
        let userShowTopic = false;
        let likesShowTopic = false;

        // CATEGORY
        const categories = topic.querySelector('div[categories]').innerHTML.split(",").map(category => category.trim());
        for (let j = 0; j < categorySelectedOptions.length; j++) {
            if (categorySelectedOptions[j] === "All categories") {
                categoryShowTopic = true;
                break;
            }
            for (let k = 0; k < categories.length; k++) {
                if (categorySelectedOptions[j] === categories[k]) {
                    categoryShowTopic = true;
                    break;
                }
            }
        }
        // USERS
        const user = topic.querySelector("a").getAttribute("author");
        for (let j = 0; j < userSelectedOptions.length; j++) {
            if (userSelectedOptions[j] === "All users") {
                userShowTopic = true;
                break;
            }
            if (userSelectedOptions[j] === user) {
                userShowTopic = true;
                break;
            }
        }
        // LIKES
        const like = topic.querySelector('.likes').getAttribute('like');
        const dislike = topic.querySelector('.dislikes').getAttribute('dislike');
        if (likesSelectedOption === "likes") {
            likesShowTopic = Boolean(parseInt(like));
        } else if (likesSelectedOption === "dislikes") {
            likesShowTopic = Boolean(parseInt(dislike));
        } else if (likesSelectedOption === "All") {
            likesShowTopic = true;
        }
        // show topic if it passes every filer
        topic.hidden = true;
        if (categoryShowTopic && userShowTopic && likesShowTopic) {
            topic.hidden = false;
        }
    }
}

function filterCategories() {
    const topics = document.getElementsByClassName("topic-group");
    const categorySelectedOptions = Array
        .from(document.getElementById("FilterCategory").options)
        .filter(option => option.selected)
        .map(option => option.value);
    if (categorySelectedOptions.length === 0) { categorySelectedOptions[0] = "All categories" }
    for (let i = 0; i < topics.length; i++) {
        const topic = topics[i];
        let categoryShowTopic = false;

        const categories = topic.querySelector('div[categories]').innerHTML.split(",").map(category => category.trim());
        for (let j = 0; j < categorySelectedOptions.length; j++) {
            if (categorySelectedOptions[j] === "All categories") {
                categoryShowTopic = true;
                break;
            }
            for (let k = 0; k < categories.length; k++) {
                if (categorySelectedOptions[j] === categories[k]) {
                    categoryShowTopic = true;
                    break;
                }
            }
        }
        topic.hidden = true;
        if (categoryShowTopic) {
            topic.hidden = false;
        }
    }
}