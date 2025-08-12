import Cookies from "js-cookie";

export function useWindowWidthCookie(cookieName = "windowWidth") {
	const width = ref(0);

	function updateWidth() {
		width.value = window.innerWidth;
		Cookies.set(cookieName, width.value.toString());
	}

	onMounted(() => {
		updateWidth();
		window.addEventListener("resize", updateWidth);
	});

	onBeforeUnmount(() => {
		window.removeEventListener("resize", updateWidth);
	});

	return { width };
}
