package main

import (
	"SchemaTreeBuilder/preparation"
	"SchemaTreeBuilder/schematree"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"runtime"
	"runtime/pprof"
	"runtime/trace"

	"github.com/spf13/cobra"
)

func main() {

	// Program initialization actions
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setup the variables where all flags will reside.
	var cpuprofile, memprofile, traceFile string // used globally
	var measureTime bool                         // used globally
	var firstNsubjects int64                     // used by build-tree
	var writeOutPropertyFreqs bool               // used by build-tree

	// Setup helper variables
	var timeCheckpoint time.Time // used globally

	// writeOutPropertyFreqs := flag.Bool("writeOutPropertyFreqs", false, "set this to write the frequency of all properties to a csv after first pass or schematree loading")

	// root command
	cmdRoot := &cobra.Command{
		Use: "SchemaTreeBuilder",

		// Execute global pre-run activities such as profiling.
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			// write cpu profile to file - open file and start profiling
			if cpuprofile != "" {
				f, err := os.Create(cpuprofile)
				if err != nil {
					log.Fatal("could not create CPU profile: ", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal("could not start CPU profile: ", err)
				}
			}

			// write trace execution to file - open file and start tracing
			if traceFile != "" {
				f, err := os.Create(traceFile)
				if err != nil {
					log.Fatal("could not create trace file: ", err)
				}
				if err := trace.Start(f); err != nil {
					log.Fatal("could not start tracing: ", err)
				}
			}

			// measure time - start measuring the time
			//   The measurements are done in such a way to not include the time for the profiles operations.
			if measureTime == true {
				timeCheckpoint = time.Now()
			}

		},

		// Close whatever profiling was running globally.
		PersistentPostRun: func(cmd *cobra.Command, args []string) {

			// measure time - stop time measurement and print the measurements
			if measureTime == true {
				fmt.Println("Execution Time:", time.Since(timeCheckpoint))
			}

			// write cpu profile to file - stop profiling
			if cpuprofile != "" {
				pprof.StopCPUProfile()
			}

			// write memory profile to file
			if memprofile != "" {
				f, err := os.Create(memprofile)
				if err != nil {
					log.Fatal("could not create memory profile: ", err)
				}
				runtime.GC() // get up-to-date statistics
				if err := pprof.WriteHeapProfile(f); err != nil {
					log.Fatal("could not write memory profile: ", err)
				}
				f.Close()
			}

			// write trace execution to file - stop tracing
			if traceFile != "" {
				trace.Stop()
			}

		},
	}

	// global flags for root command
	cmdRoot.PersistentFlags().StringVar(&cpuprofile, "cpuprofile", "", "write cpu profile to `file`")
	cmdRoot.PersistentFlags().StringVar(&memprofile, "memprofile", "", "write memory profile to `file`")
	cmdRoot.PersistentFlags().StringVar(&traceFile, "trace", "", "write execution trace to `file`")
	cmdRoot.PersistentFlags().BoolVarP(&measureTime, "time", "t", false, "measure time of command execution")

	// subcommand build-tree
	cmdBuildTreeTyped := &cobra.Command{
		Use:   "build-tree-typed <dataset>",
		Short: "Build the SchemaTree model with types",
		Long: "A SchemaTree model will be built using the file provided in <dataset>." +
			" The dataset should be a N-Triple of Items.\nTwo output files will be" +
			" generated in the same directory as <dataset> and with suffixed names, namely:" +
			" '<dataset>.firstPass.bin' and '<dataset>.schemaTree.typed.bin'",
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			inputDataset := &args[0]

			// Create the tree output file by using the input dataset.
			schema, err := schematree.Create(*inputDataset, uint64(firstNsubjects), true, 0)
			if err != nil {
				log.Panicln(err)
			}

			if writeOutPropertyFreqs {
				propFreqsPath := *inputDataset + ".propertyFreqs.csv"
				schema.WritePropFreqs(propFreqsPath)
				fmt.Printf("Wrote PropertyFreqs to %s\n", propFreqsPath)

				typeFreqsPath := *inputDataset + ".typeFreqs.csv"
				schema.WriteTypeFreqs(typeFreqsPath)
				fmt.Printf("Wrote PropertyFreqs to %s\n", typeFreqsPath)
			}

		},
	}

	// cmdBuildTree.Flags().StringVarP(&inputDataset, "dataset", "d", "", "`path` to the dataset file to parse")
	// cmdBuildTree.MarkFlagRequired("dataset")
	cmdBuildTreeTyped.Flags().Int64VarP(&firstNsubjects, "first", "n", 0, "only parse the first `n` subjects") // TODO: handle negative inputs
	cmdBuildTreeTyped.Flags().BoolVarP(
		&writeOutPropertyFreqs, "write-frequencies", "f", false,
		"write all property frequencies to a csv file named '<dataset>.propertyFreqs.csv' after the SchemaTree is built",
	)
	cmdBuildTreeTyped.Flags().IntVar(&numPointersInNode, "num-pointers-typed", 3, "The number of pointers sotred directly in the node") // TODO: handle negative inputs

	// subcommand build-tree
	cmdBuildTree := &cobra.Command{
		Use:   "build-tree <dataset>",
		Short: "Build the SchemaTree model",
		Long: "A SchemaTree model will be built using the file provided in <dataset>." +
			" The dataset should be a N-Triple of Items.\nTwo output files will be" +
			" generated in the same directory as <dataset> and with suffixed names, namely:" +
			" '<dataset>.firstPass.bin' and '<dataset>.schemaTree.bin'",
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			inputDataset := &args[0]

			// Create the tree output file by using the input dataset.
			schema, err := schematree.Create(*inputDataset, uint64(firstNsubjects), false, 0)
			if err != nil {
				log.Panicln(err)
			}

			if writeOutPropertyFreqs {
				propFreqsPath := *inputDataset + ".propertyFreqs.csv"
				schema.WritePropFreqs(propFreqsPath)
				fmt.Printf("Wrote PropertyFreqs to %s\n", propFreqsPath)
			}

		},
	}
	// cmdBuildTree.Flags().StringVarP(&inputDataset, "dataset", "d", "", "`path` to the dataset file to parse")
	// cmdBuildTree.MarkFlagRequired("dataset")
	cmdBuildTree.Flags().Int64VarP(&firstNsubjects, "first", "n", 0, "only parse the first `n` subjects") // TODO: handle negative inputs
	cmdBuildTree.Flags().BoolVarP(
		&writeOutPropertyFreqs, "write-frequencies", "f", false,
		"write all property frequencies to a csv file named '<dataset>.propertyFreqs.csv' after the SchemaTree is built",
	)

	// subcommand split-dataset
	cmdSplitDataset := &cobra.Command{
		Use:   "split-dataset",
		Short: "Split a dataset using various methods",
		Long: "Select the method with which to split a N-Triple dataset file and" +
			" generate multiple smaller datasets in the same directory and with" +
			" suffixed names. Suffixes depend on chosen splitter method.",
		Args: cobra.NoArgs,
	}

	// subsubcommand split-dataset by-prefix
	cmdSplitDatasetByPrefix := &cobra.Command{
		Use:   "by-prefix <dataset>",
		Short: "Split a dataset according to the prefix of the subject",
		Long: "Split a N-Triple <dataset> file into three files according to the preset of the subject" +
			" into: item, prop and misc.\nThe split files are generated in the same directory" +
			" as the <dataset>, stripped of their compression extension and given the following" +
			" names: <base>-item.nt.gz, <base>-prop.nt.gz, <dataset>-misc.nt.gz",
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			inputDataset := &args[0]

			// Make the split
			sStats, err := preparation.SplitByPrefix(*inputDataset)
			if err != nil {
				log.Panicln(err)
			}

			// Prepare and output the stats for it
			totalCount := float64(sStats.ItemCount + sStats.PropCount + sStats.MiscCount)
			fmt.Println("Split dataset by prefix:")
			fmt.Printf("  item: %d (%f)\n", sStats.ItemCount, float64(sStats.ItemCount)/totalCount)
			fmt.Printf("  prop: %d (%f)\n", sStats.PropCount, float64(sStats.PropCount)/totalCount)
			fmt.Printf("  misc: %d (%f)\n", sStats.MiscCount, float64(sStats.MiscCount)/totalCount)

		},
	}

	// subsubcommand split-dataset 1-in-n
	cmdSplitDatasetBySampling := &cobra.Command{
		Use:   "1-in-n <dataset>",
		Short: "Split a dataset using systematic sampling",
		Long: "Split a N-Triple <dataset> file into two files where every Nth subject goes into" +
			" one file and the rest into the second file.\nThe split files are generated in the same directory" +
			" as the <dataset>, stripped of their compression extension and given the following" +
			" names: <base>-1in<n>-test.nt.gz, <base>-1in<n>-train.nt.gz\n" +
			"This method assumes that all entries for a given subject are defined in contiguous lines.",
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			inputDataset := &args[0]
			preparation.SplitBySampling(*inputDataset, int64(everyNthSubject))
		},
	}
	cmdSplitDatasetBySampling.Flags().UintVarP(&everyNthSubject, "nth", "n", 1000, "split every N-th subject")

	// subcommand filter-dataset
	cmdFilterDataset := &cobra.Command{
		Use:   "filter-dataset",
		Short: "Filter a dataset using various methods",
		Long:  "Filter the dataset for the purpose of building other models.",
		Args:  cobra.NoArgs,
	}

	// subsubcommand filter-dataset for-schematree
	cmdFilterDatasetForSchematree := &cobra.Command{
		Use:   "for-schematree <dataset>",
		Short: "Prepare the dataset for inclusion in the SchemaTree",
		Long: "Remove entries that should not be considered by the SchemaTree builder.\nThe new file is" +
			" generated in the same directory as the <dataset>, stripped of their compression extension" +
			" and given the following name: <base>-filtered.nt.gz",
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			inputDataset := &args[0]

			// Execute the filter
			sStats, err := preparation.FilterForSchematree(*inputDataset)
			if err != nil {
				log.Panicln(err)
			}

			// Prepare and output the stats for it
			totalCount := float64(sStats.KeptCount + sStats.LostCount)
			fmt.Println("Filter dataset for schematree:")
			fmt.Printf("  kept: %d (%f)\n", sStats.KeptCount, float64(sStats.KeptCount)/totalCount)
			fmt.Printf("  lost: %d (%f)\n", sStats.LostCount, float64(sStats.LostCount)/totalCount)

		},
	}

	// putting the command hierarchy together
	cmdRoot.AddCommand(cmdSplitDataset)
	cmdSplitDataset.AddCommand(cmdSplitDatasetByPrefix)
	cmdSplitDataset.AddCommand(cmdSplitDatasetBySampling)
	cmdRoot.AddCommand(cmdFilterDataset)
	cmdFilterDataset.AddCommand(cmdFilterDatasetForSchematree)
	cmdRoot.AddCommand(cmdBuildTree)
	cmdRoot.AddCommand(cmdBuildTreeTyped)

	// Start the CLI application
	cmdRoot.Execute()

}

func waitForReturn() {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	sentence, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(sentence))
	}
}
