<script setup lang="ts">
const messages = ref<string[]>([]);
const article = ref<string>(`
# Install Arch + niri + btrfs + refind

## 1. 准备网络

- 有线：\`ip link\` 检查/连接网线

- 无线（使用 \`iwctl\` / iwd）
  - 列设备：\`device list\`
  - 打开设备/适配器：\`device <name> set-property Powered on\` / \`adapter <name> set-property Powered on\`
  - 扫描：\`station <name> scan\`
  - 列出网络：\`station <name> get-networks\`
  - 连接：\`station <name> connect <SSID>\`

## 2. 更换 Pacman 国内镜像（提高速度）

- 安装 reflector：\`pacman -Sy reflector\`

- 生成镜像列表：

  \`\`\`bash
  reflector --sort rate --threads 10 -c China --save /etc/pacman.d/mirrorlist
  \`\`\`

## 3. 查看磁盘并分区

- 查看：\`fdisk -l\`（或用 \`lsblk\` 等工具）用于确认磁盘与分区

- 使用 \`gdisk\` 分区

  1. **启动 \`gdisk\`**

     \`\`\`bash
     gdisk /dev/sdX  # sdX 是要分区的磁盘（比如 /dev/nvme0n1）
     \`\`\`

  2. **查看现有分区**

     在 \`gdisk\` 提示符下输入 \`p\`，它会列出当前磁盘的分区表

  3. **删除旧分区**（如果需要）

     如果磁盘上有旧分区，输入 \`d\` 来删除，选择要删除的分区号

  4. **创建新分区**

     输入 \`n\` 来创建新分区。你需要选择分区号、起始和结束位置

  5. **设置分区类型**

     使用 \`t\` 来设置分区类型。常用的类型有：

     - \`EF00\`：EFI 系统分区（ESP）(300MB)

     - \`8300\`：Linux 文件系统
     - \`8200\`：Linux Swap 分区（如果需要交换分区）

  6. **保存并退出**

     输入 \`w\` 来写入分区表并退出 \`gdisk\`

## 4. 格式化分区

使用 \`mkfs\` 命令来格式化分区

- **根分区**

  \`\`\`bash
  mkfs.btrfs /dev/nvme0n1p2  # 假设根分区是 /dev/nvme0n1p2
  \`\`\`

- **EFI 分区**（使用 UEFI 启动）：

  \`\`\`bash
  mkfs.fat -F32 /dev/nvme0n1p1  # 假设 EFI 分区是 /dev/nvme0n1p1
  \`\`\`

## 5. 创建子卷

\`\`\`bash
mount /dev/nvme0n1p2 /mnt
btrfs subvolume create /mnt/@
btrfs subvolume create /mnt/@home
btrfs subvolume create /mnt/@log
btrfs subvolume create /mnt/@pkg
umount /mnt
\`\`\`

## 6. 挂载文件系统

\`\`\`bash
mount -o noatime,compress=zstd,space_cache=v2,subvol=@ /dev/nvme0n1p2 /mnt
mkdir -p /mnt/{boot/efi,home,var/log,var/cache/pacman/pkg}
mount -o noatime,compress=zstd,space_cache=v2,subvol=@home /dev/nvme0n1p2 /mnt/home
mount -o noatime,compress=zstd,space_cache=v2,subvol=@log /dev/nvme0n1p2 /mnt/var/log
mount -o noatime,compress=zstd,space_cache=v2,subvol=@pkg /dev/nvme0n1p2 /mnt/var/cache/pacman/pkg
\`\`\`

## 7. 基本安装（示例）

\`\`\`bash  
pacstrap -K /mnt base linux linux-firmware btrfs-progs
\`\`\`

可选

- \`intel-ucode\` or \`amd-ucode\`

- neovim
- base-devel
- nvidia-dkms
- egl-wayland
- nvidia-utils

then

\`\`\`bash
genfstab -U /mnt > /mnt/etc/fstab
\`\`\`

`); // 用户输入的文章

const summarize = async () => {
  const data = await $fetch<{ sessionId: string }>("/api/v1/model/summarize", {
    method: "post",
    body: {
      content: article.value,
    },
  });

  const es = new EventSource("/api/v1/model/summarize/" + data.sessionId);

  es.onmessage = (event) => {
    setTimeout(() => {
      messages.value.push(event.data);
      console.log(event.data);
    }, 1000);
  };

  es.addEventListener("done", () => {
    console.log("生成完成");
    es.close();
  });

  es.addEventListener("error", (e) => {
    console.error("SSE 错误:", e);
    es.close();
  });
};
</script>

<template>
  <div class="container">
    <h1>实时消息</h1>
    <TransitionGroup name="msgs">
      <span v-for="msg in messages" :key="msg">{{ msg }}</span>
    </TransitionGroup>
    <button @click="summarize">generate</button>
  </div>
</template>

<style lang="less" scoped>
.container {
  width: 100vw;
  height: 100vh;
}

.msgs-enter-active,
.msgs-leave-active {
  transition: all 0.5s ease;
}
.msgs-enter-from,
.msgs-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
