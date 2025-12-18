export default defineNuxtPlugin(() => {
  const { done } = useLoadingBar();

  useRuntimeHook("page:finish", () => {
    setTimeout(() => {
      done();
    }, 500);
  });
});
