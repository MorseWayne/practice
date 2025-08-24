import yt_dlp
import os
import sys


def download_youtube_video(
    url,
    output_dir="./downloads",
    format=(
        "bestvideo[height<=1080][fps<=30][ext=mp4]+bestaudio[ext=m4a]/"  # 首选：1080p 30帧MP4视频+最佳音频
        "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/"  # 备选：1080p MP4视频+最佳音频
        "best[height<=1080][ext=mp4]/best"
    ),
    concurrent_fragments=8,  # 并行下载的片段数量
):  # 备选：1080p MP4格式或最佳可用格式
    """
    从YouTube下载视频，保存为指定质量的MP4格式

    参数:
        url (str): YouTube视频URL
        output_dir (str): 输出目录，默认为./downloads
        format (str): 视频格式规范，默认为1080p 30帧的MP4格式
        concurrent_fragments (int): 并行下载的片段数量，默认为8，可提高下载速度
    """
    # 创建输出目录（如果不存在）
    os.makedirs(output_dir, exist_ok=True)

    # 设置下载选项
    ydl_opts = {
        "format": format,
        "outtmpl": os.path.join(output_dir, "%(title)s.%(ext)s"),
        "merge_output_format": "mp4",
        "postprocessors": [
            {
                "key": "FFmpegVideoConvertor",
                "preferedformat": "mp4",
            }
        ],
        # 配置并行下载以提高速度
        "concurrent_fragment_downloads": concurrent_fragments,  # 启用多片段并行下载
        "allow_multiple_video_streams": True,  # 允许多视频流
        "allow_multiple_audio_streams": True,  # 允许多音频流
        # 以下选项可进一步优化下载性能
        "socket_timeout": 15,  # 套接字超时时间(秒)
        "retries": 10,  # 下载失败重试次数
        "retry_sleep": 1,  # 重试间隔(秒)
    }

    try:
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            print(f"开始下载视频: {url}")
            ydl.download([url])
            print(f"视频已成功下载到: {output_dir}")
            return True
    except Exception as e:
        print(f"下载失败: {str(e)}")
        return False


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("用法: python download_video.py <YouTube视频URL>")
        print(
            "示例: python download_video.py https://www.youtube.com/watch?v=dQw4w9WgXcQ"
        )
        sys.exit(1)

    video_url = sys.argv[1]
    download_youtube_video(video_url)

# 注意: 运行此脚本前需要安装yt-dlp和ffmpeg
# 安装命令: pip install yt-dlp
# 对于ffmpeg，请根据您的操作系统安装相应版本: https://ffmpeg.org/download.html
