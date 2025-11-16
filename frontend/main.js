document.addEventListener("DOMContentLoaded", () => {
  const API_URL = "http://127.0.0.1:8080";

  // ====== Хранение токена (после логина) ======
  function getToken() {
    return localStorage.getItem("token");
  }

  function setToken(token) {
    if (!token) return;
    localStorage.setItem("token", token);
  }

  // ====== Базовый помощник для запросов ======
  async function apiRequest(
    path,
    { method = "GET", params, body, auth = false } = {}
  ) {
    let url = API_URL + path;

    if (params && method === "GET") {
      const qs = new URLSearchParams(params).toString();
      if (qs) url += "?" + qs;
    }

    const headers = {};

    if (body && method !== "GET") {
      headers["Content-Type"] = "application/json";
    }

    if (auth) {
      const token = getToken();
      if (token) {
        headers["Authorization"] = "Bearer " + token;
      }
    }

    const res = await fetch(url, {
      method,
      headers,
      body: body && method !== "GET" ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      const text = await res.text().catch(() => "");
      throw new Error(`Ошибка ${res.status}: ${text || res.statusText}`);
    }

    try {
      return await res.json();
    } catch {
      return null;
    }
  }

  // ====== МОДАЛКИ ======
  const loginModal = document.getElementById("login-modal");
  const ngoModal = document.getElementById("ngo-modal");
  const registerModal = document.getElementById("register-modal");
  const profileModal = document.getElementById("profile-modal");

  const openLoginBtn = document.getElementById("open-login-modal");
  const openNgoBtn = document.getElementById("open-ngo-modal");
  const openRegisterFromLogin = document.getElementById(
    "open-register-from-login"
  );
  const openProfileBtn = document.getElementById("open-profile-modal");

  function openModal(modal) {
    if (modal) modal.classList.add("modal--open");
  }

  function closeModal(modal) {
    if (modal) modal.classList.remove("modal--open");
  }

  openLoginBtn?.addEventListener("click", () => openModal(loginModal));
  openNgoBtn?.addEventListener("click", () => openModal(ngoModal));

  openRegisterFromLogin?.addEventListener("click", () => {
    closeModal(loginModal);
    openModal(registerModal);
  });

  // Открытие профиля
  openProfileBtn?.addEventListener("click", async () => {
    const token = getToken();
    if (!token) {
      alert("Сначала войдите в систему, чтобы открыть профиль.");
      openModal(loginModal);
      return;
    }

    try {
      const data = await apiRequest("/me", {
        method: "GET",
        auth: true,
      });

      // ожидаем от бэка поля: email, fullName, ngoName
      const emailEl = document.getElementById("profile-email");
      const fullnameEl = document.getElementById("profile-fullname");
      const ngoEl = document.getElementById("profile-ngo");

      if (emailEl) emailEl.textContent = data?.email || "—";
      if (fullnameEl) fullnameEl.textContent = data?.fullName || "—";
      if (ngoEl) ngoEl.textContent = data?.ngoName || "НКО не привязана";

      openModal(profileModal);
    } catch (err) {
      console.error(err);
      alert(
        "Не удалось загрузить профиль. Возможно, сессия истекла — войдите заново."
      );
      openModal(loginModal);
    }
  });

  document.querySelectorAll("[data-close-modal]").forEach((el) => {
    el.addEventListener("click", () => {
      const modal = el.closest(".modal");
      closeModal(modal);
    });
  });

  document.querySelectorAll(".modal__backdrop").forEach((backdrop) => {
    backdrop.addEventListener("click", () => {
      const modal = backdrop.closest(".modal");
      closeModal(modal);
    });
  });

  document.addEventListener("keydown", (e) => {
    if (e.key === "Escape") {
      [loginModal, ngoModal, registerModal, profileModal].forEach(closeModal);
    }
  });

  // ====== РЕНДЕР СПИСКА НКО (ответ на /ngos) ======
  const ngoListEl = document.querySelector(".ngo-list");

  function renderNgoList(ngos) {
    if (!ngoListEl) return;
    ngoListEl.innerHTML = "";

    if (!Array.isArray(ngos) || ngos.length === 0) {
      ngoListEl.innerHTML =
        '<p style="font-size:14px;color:#666;">Ничего не найдено</p>';
      return;
    }

    ngos.forEach((ngo) => {
      const card = document.createElement("div");
      card.className = "ngo-card";

      const name = ngo.name || "Без названия";
      const category = ngo.category || "Не указано";
      const desc =
        ngo.description || ngo.volunteerDescription || "Описание отсутствует";

      card.innerHTML = `
        <h3>${name}</h3>
        <p class="cat">Категория: ${category}</p>
        <p class="desc">${desc}</p>
        <button class="btn-small" type="button">Показать на карте</button>
      `;

      ngoListEl.appendChild(card);
    });
  }

  // ====== ФОРМА ПОИСКА / ФИЛЬТРОВ (/ngos) ======
  const searchForm = document.getElementById("search-form");

  searchForm?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(searchForm);
    const params = Object.fromEntries(formData.entries());

    try {
      const data = await apiRequest("/ngos", {
        method: "GET",
        params,
      });

      console.log("Ответ /ngos:", data);
      renderNgoList(data);
    } catch (err) {
      console.error(err);
      alert("Ошибка при загрузке списка НКО");
    }
  });

  // Можно сразу подгрузить все НКО без фильтров при загрузке страницы:
  (async () => {
    try {
      const data = await apiRequest("/ngos", { method: "GET" });
      renderNgoList(data);
    } catch (e) {
      console.warn("Не удалось загрузить НКО по умолчанию:", e.message);
    }
  })();

  // ====== ФОРМА ЛОГИНА (/login) ======
  const loginForm = document.getElementById("login-form");

  loginForm?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(loginForm);
    const body = Object.fromEntries(formData.entries());

    try {
      const data = await apiRequest("/login", {
        method: "POST",
        body,
      });

      console.log("Ответ /login:", data);

      if (data && data.token) {
        setToken(data.token);
      }

      alert("Успешный вход");
      closeModal(loginModal);
    } catch (err) {
      console.error(err);
      alert("Ошибка входа. Проверьте данные.");
    }
  });

  // ====== ФОРМА РЕГИСТРАЦИИ (/register) ======
  const registerForm = document.getElementById("register-form");

  registerForm?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(registerForm);
    const body = Object.fromEntries(formData.entries());

    try {
      const data = await apiRequest("/register", {
        method: "POST",
        body,
      });

      console.log("Ответ /register:", data);
      alert("Аккаунт создан");

      closeModal(registerModal);
      openModal(loginModal);
    } catch (err) {
      console.error(err);
      alert("Ошибка регистрации");
    }
  });

  // ====== ФОРМА НКО (/ngo) ======
  const ngoForm = document.getElementById("ngo-form");

  ngoForm?.addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = new FormData(ngoForm);
    const body = Object.fromEntries(formData.entries());

    try {
      const data = await apiRequest("/ngo", {
        method: "POST",
        body,
        auth: true,
      });

      console.log("Ответ /ngo:", data);
      alert("НКО сохранена");

      closeModal(ngoModal);

      // Обновляем список НКО
      try {
        const ngos = await apiRequest("/ngos", { method: "GET" });
        renderNgoList(ngos);
      } catch (e2) {
        console.warn("Не удалось обновить список НКО:", e2.message);
      }
    } catch (err) {
      console.error(err);
      alert("Ошибка при сохранении НКО");
    }
  });
});
