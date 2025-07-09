package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new project from this template",
	Long: `Initialize a new project from this template by replacing the module name
and updating all import paths throughout the codebase.

Examples:
  gin-starter init my-awesome-app
  gin-starter init github.com/myorg/my-app
  gin-starter init gitlab.com/company/project`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

type ProjectConfig struct {
	ModuleName    string
	ProjectName   string
	TargetDir     string
	TemplateDir   string
	ReplaceModule string
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("output", "o", "", "Output directory (default: project name)")
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite existing directory")
	initCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
}

func runInit(cmd *cobra.Command, args []string) error {
	moduleName := args[0]
	
	// 获取标志
	outputDir, _ := cmd.Flags().GetString("output")
	force, _ := cmd.Flags().GetBool("force")
	interactive, _ := cmd.Flags().GetBool("interactive")
	
	// 解析项目名
	projectName := extractProjectName(moduleName)
	if outputDir == "" {
		outputDir = projectName
	}
	
	// 获取当前目录作为模版目录
	templateDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	
	config := &ProjectConfig{
		ModuleName:    moduleName,
		ProjectName:   projectName,
		TargetDir:     outputDir,
		TemplateDir:   templateDir,
		ReplaceModule: "github.com/iswangwenbin/gin-starter",
	}
	
	// 交互模式
	if interactive {
		if err := runInteractiveMode(config); err != nil {
			return err
		}
	}
	
	// 验证配置
	if err := validateConfig(config); err != nil {
		return err
	}
	
	// 检查目标目录
	if err := checkTargetDirectory(config.TargetDir, force); err != nil {
		return err
	}
	
	// 执行初始化
	fmt.Printf("🚀 Initializing new project '%s'...\n", config.ProjectName)
	fmt.Printf("📂 Module name: %s\n", config.ModuleName)
	fmt.Printf("📁 Target directory: %s\n", config.TargetDir)
	fmt.Println()
	
	if err := initializeProject(config); err != nil {
		return fmt.Errorf("failed to initialize project: %w", err)
	}
	
	fmt.Println("✅ Project initialized successfully!")
	fmt.Printf("📋 Next steps:\n")
	fmt.Printf("   cd %s\n", config.TargetDir)
	fmt.Printf("   make deps\n")
	fmt.Printf("   make run\n")
	
	return nil
}

func extractProjectName(moduleName string) string {
	// 从模块名提取项目名
	parts := strings.Split(moduleName, "/")
	return parts[len(parts)-1]
}

func runInteractiveMode(config *ProjectConfig) error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("🎯 Interactive Project Setup")
	fmt.Println("===============================")
	
	// 确认模块名
	fmt.Printf("Module name [%s]: ", config.ModuleName)
	if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
		config.ModuleName = strings.TrimSpace(input)
		config.ProjectName = extractProjectName(config.ModuleName)
	}
	
	// 确认项目名
	fmt.Printf("Project name [%s]: ", config.ProjectName)
	if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
		config.ProjectName = strings.TrimSpace(input)
	}
	
	// 确认目标目录
	fmt.Printf("Target directory [%s]: ", config.TargetDir)
	if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
		config.TargetDir = strings.TrimSpace(input)
	}
	
	return nil
}

func validateConfig(config *ProjectConfig) error {
	if config.ModuleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	
	if config.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	
	if config.TargetDir == "" {
		return fmt.Errorf("target directory cannot be empty")
	}
	
	// 验证模块名格式
	if !isValidModuleName(config.ModuleName) {
		return fmt.Errorf("invalid module name format: %s", config.ModuleName)
	}
	
	return nil
}

func isValidModuleName(name string) bool {
	// 简单的模块名验证
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._/-]+$`, name)
	return matched && !strings.Contains(name, " ")
}

func checkTargetDirectory(targetDir string, force bool) error {
	if _, err := os.Stat(targetDir); err == nil {
		if !force {
			return fmt.Errorf("directory '%s' already exists, use --force to overwrite", targetDir)
		}
		fmt.Printf("⚠️  Directory '%s' exists, will overwrite...\n", targetDir)
	}
	return nil
}

func initializeProject(config *ProjectConfig) error {
	// 创建目标目录
	if err := os.MkdirAll(config.TargetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	
	// 复制文件并替换内容
	return copyAndReplaceFiles(config)
}

func copyAndReplaceFiles(config *ProjectConfig) error {
	// 需要排除的文件和目录
	excludePatterns := []string{
		".git",
		".idea",
		".vscode",
		"*.log",
		"logs",
		"tmp",
		"bin",
		"coverage.out",
		"coverage.html",
		".DS_Store",
		"Thumbs.db",
		config.TargetDir, // 避免复制到自己
	}
	
	return filepath.Walk(config.TemplateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 计算相对路径
		relPath, err := filepath.Rel(config.TemplateDir, path)
		if err != nil {
			return err
		}
		
		// 跳过排除的文件
		if shouldExclude(relPath, excludePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		// 目标路径
		targetPath := filepath.Join(config.TargetDir, relPath)
		
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		
		return copyAndReplaceFile(path, targetPath, config)
	})
}

func shouldExclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func copyAndReplaceFile(srcPath, dstPath string, config *ProjectConfig) error {
	// 读取源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	// 创建目标文件
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}
	
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	// 获取文件信息
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	
	// 如果是文本文件，替换内容
	if isTextFile(srcPath) {
		content, err := io.ReadAll(srcFile)
		if err != nil {
			return err
		}
		
		// 替换模块名
		newContent := string(content)
		newContent = strings.ReplaceAll(newContent, config.ReplaceModule, config.ModuleName)
		
		// 替换项目名称（在某些配置文件中）
		newContent = strings.ReplaceAll(newContent, "gin-starter", config.ProjectName)
		
		// 写入新内容
		if _, err := dstFile.WriteString(newContent); err != nil {
			return err
		}
	} else {
		// 二进制文件直接复制
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}
	}
	
	// 设置文件权限
	return os.Chmod(dstPath, srcInfo.Mode())
}

func isTextFile(filename string) bool {
	textExtensions := []string{
		".go", ".md", ".txt", ".yml", ".yaml", ".json", ".toml",
		".sql", ".sh", ".bat", ".dockerfile", ".gitignore",
		".gitattributes", ".editorconfig", ".golangci.yml",
		".env", ".env.example", "Makefile", "README",
	}
	
	ext := strings.ToLower(filepath.Ext(filename))
	base := strings.ToLower(filepath.Base(filename))
	
	for _, textExt := range textExtensions {
		if ext == textExt || base == strings.TrimPrefix(textExt, ".") {
			return true
		}
	}
	
	return false
}